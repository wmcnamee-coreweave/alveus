// notes:
// https://github.com/argoproj/argo-cd/blob/04794332d28e44c18c89b83e9c54f0bcd69bcc43/cmd/argocd/commands/app.go#L1985
// https://github.com/marketplace/actions/manual-workflow-approval

package river

import (
	"fmt"

	"github.com/ghostsquad/alveus/api/v1alpha1"
	"github.com/ghostsquad/alveus/internal/integrations/github"
	"github.com/ghostsquad/alveus/internal/river/stages"
	"github.com/ghostsquad/alveus/internal/util"
)

type River struct {
	wf             *github.Workflow
	childWorkflows []*github.Workflow
}

func (r *River) GetWorkflow() *github.Workflow {
	return r.wf
}

func (r *River) GetChildWorkflows() []*github.Workflow {
	return r.childWorkflows
}

func (r *River) GetAllWorkflows() []*github.Workflow {
	allWorkflows := []*github.Workflow{r.wf}
	return append(allWorkflows, r.childWorkflows...)
}

func New(service v1alpha1.Service) (*River, error) {
	r := River{}

	wf := github.Workflow{
		Name: service.Name,
		On: github.WorkflowOn{
			Call: github.OnCall{
				Inputs:  map[string]github.CallInput{},
				Outputs: map[string]github.CallOutput{},
			},
			Dispatch: github.OnDispatch{
				Inputs: map[string]github.OnDispatchInput{},
			},
		},
		Concurrency: github.Concurrency{},
		Defaults:    github.Defaults{},
		Jobs: map[string]github.Job{
			"source": {
				Permissions: github.Permissions{},
				RunsOn:      nil,
				Concurrency: github.Concurrency{},
				Outputs:     nil,
				Env:         nil,
				Defaults:    github.Defaults{},
				Steps: []github.Step{
					{
						Uses: "actions/checkout@v4",
						With: map[string]any{
							"fetch-depth": 1,
						},
					},
					{
						Run: "alveus generate applications",
					},
				},
			},
		},
	}

	var previousWorkflowName string

	for _, g := range service.DestinationGroups {
		gwf := github.Workflow{
			Name: util.Join("-", service.Name, g.Name),
			On: github.WorkflowOn{
				Call: github.OnCall{
					Inputs:  map[string]github.CallInput{},
					Outputs: map[string]github.CallOutput{},
				},
			},
			Concurrency: github.Concurrency{},
			Defaults:    github.Defaults{},
			Jobs:        map[string]github.Job{},
		}

		var needs []string
		if previousWorkflowName != "" {
			needs = append(needs, previousWorkflowName)
		}

		wf.Jobs[gwf.Name] = github.Job{
			Name: gwf.Name,
			Uses: gwf.GetFullFilename(),
		}

		previousWorkflowName = gwf.Name

		for _, d := range g.Destinations {
			stg, err := stages.New(gwf.Name, d)
			if err != nil {
				return nil, fmt.Errorf("creating stage: %w", err)
			}
			swf := stg.GetWorkflow()
			// https://docs.github.com/en/actions/how-tos/sharing-automations/reusing-workflows
			gwf.Jobs[swf.Name] = github.Job{
				Name: swf.Name,
				Uses: swf.GetFullFilename(),
				With: map[string]any{},
			}

			r.childWorkflows = append(r.childWorkflows, swf)
		}

		r.childWorkflows = append(r.childWorkflows, &gwf)
	}

	r.wf = &wf

	return &r, nil
}
