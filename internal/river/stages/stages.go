package stages

import (
	"fmt"

	"github.com/ghostsquad/alveus/api/v1alpha1"
	"github.com/ghostsquad/alveus/internal/constants"
	"github.com/ghostsquad/alveus/internal/integrations/github"
	"github.com/ghostsquad/alveus/internal/river/jobs"
	"github.com/ghostsquad/alveus/internal/util"
)

const (
	StageInputTargetRevision = "target-revision"
)

// Stage is an abstraction around a GitHub workflow
// the workflow updates the associated ArgoCD Application yaml file in the /state directory
// the same changes are also pushed to ArgoCD to trigger the application to sync
type Stage struct {
	filename string
	wf       *github.Workflow
}

func (s *Stage) GetWorkflow() *github.Workflow {
	return s.wf
}

func (s *Stage) GetFilename() string {
	return s.filename
}

type StageOptions struct {
	ExcludeCommitJob bool
}

type StageOption func(*StageOptions)

func New(parentName string, destination v1alpha1.Destination, options ...StageOption) (Stage, error) {
	opts := &StageOptions{}
	for _, o := range options {
		o(opts)
	}

	if !opts.ExcludeCommitJob {
	}

	commitJob, err := jobs.NewCommitChanges(jobs.CommitChangesInput{
		JobName:       "",
		CommitMessage: "",
		Ref:           "",
	}, jobs.CommitChangesWithIntermediateSteps(github.Step{
		Run: constants.CLIName + " update application",
	}))

	if err != nil {
		return Stage{}, fmt.Errorf("creating job: %w", err)
	}

	wf := github.Workflow{
		Name: util.Join("-", parentName, destination.String()),
		On: github.WorkflowOn{
			Call:     github.OnCall{},
			Dispatch: github.OnDispatch{},
		},
		Concurrency: github.Concurrency{},
		Defaults:    github.Defaults{},
		Jobs: map[string]github.Job{
			"-": commitJob,
		},
	}

	return Stage{
		wf:       &wf,
		filename: wf.GetFullFilename(),
	}, nil
}
