package jobs

import (
	"errors"
	"fmt"

	"github.com/ghostsquad/alveus/internal/integrations/github"
	"github.com/ghostsquad/alveus/internal/integrations/github/expressions"
	"github.com/ghostsquad/alveus/internal/util"
)

type CommitChangesOptions struct {
	IntermediateSteps []github.Step
	Permissions       github.Permissions
	Needs             []string
	If                string
	RunsOn            []string
	Outputs           map[string]string
	Environment       github.Environment
	Concurrency       github.Concurrency
	Env               map[string]string
	Defaults          github.Defaults
	PostSteps         []github.Step
}

type CommitChangeOption func(*CommitChangesOptions)

func CommitChangesWithIntermediateSteps(intermediateSteps ...github.Step) CommitChangeOption {
	return func(options *CommitChangesOptions) {
		options.IntermediateSteps = append(options.IntermediateSteps, intermediateSteps...)
	}
}

type CommitChangesInput struct {
	JobName       string
	CommitMessage string
	Ref           string
}

func (input *CommitChangesInput) Validate() error {
	var errs []error

	if input.JobName == "" {
		errs = append(errs, fmt.Errorf("job name is required"))
	}

	if input.Ref == "" {
		errs = append(errs, fmt.Errorf("ref is required"))
	}

	if input.CommitMessage == "" {
		errs = append(errs, fmt.Errorf("commit message is required"))
	}

	return errors.Join(errs...)
}

func NewCommitChanges(input CommitChangesInput, options ...CommitChangeOption) (github.Job, error) {
	opts := &CommitChangesOptions{
		Outputs: make(map[string]string),
		Env:     make(map[string]string),
	}
	for _, o := range options {
		o(opts)
	}

	if err := input.Validate(); err != nil {
		return github.Job{}, fmt.Errorf("invalid input: %w", err)
	}

	steps := []github.Step{
		{
			Uses: "actions/checkout@v4",
			With: map[string]any{
				"fetch-depth": 1,
				"ssh-key":     expressions.Secrets("WRITE_DEPLOY_KEY").String(),
				"ref":         input.Ref,
			},
		},
		{
			Name: "configure-git",
			Run: util.Join("\n",
				"git config --local user.email 'github-actions@github.com'",
				"git config --local user.name 'GitHub Actions'",
			),
		},
	}

	steps = append(steps, opts.IntermediateSteps...)

	steps = append(steps, github.Step{
		Run: util.Join("\n",
			"git add .",
			fmt.Sprintf("git commit -m %q", input.CommitMessage),
			"git push",
		),
	})

	steps = append(steps, opts.PostSteps...)

	job := github.Job{
		Name:        input.JobName,
		Permissions: opts.Permissions,
		Needs:       opts.Needs,
		If:          opts.If,
		RunsOn:      opts.RunsOn,
		Environment: opts.Environment,
		Concurrency: opts.Concurrency,
		Outputs:     opts.Outputs,
		Env:         opts.Env,
		Defaults:    opts.Defaults,
		Steps:       steps,
	}

	return job, nil
}
