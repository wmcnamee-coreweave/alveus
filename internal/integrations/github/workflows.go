package github

import (
	"github.com/cakehappens/gocto"

	"github.com/ghostsquad/alveus/api/v1alpha1"
)

func NewWorkflows(service v1alpha1.Service) []gocto.Workflow {
	var workflows []gocto.Workflow

	top := gocto.Workflow{
		Name: service.Name,
		On: gocto.WorkflowOn{
			Dispatch: gocto.OnDispatch{},
			Push: gocto.OnPush{
				OnPaths:    gocto.OnPaths{},
				OnBranches: gocto.OnBranches{},
				OnTags:     gocto.OnTags{},
			},
		},
		Jobs: make(map[string]gocto.Job),
	}

	var prevGroupJob *gocto.Job
	for _, dg := range service.DestinationGroups {
		dgWf, subWfs := newDeploymentGroupWorkflows(dg)
		workflows = append(workflows, dgWf)
		workflows = append(workflows, subWfs...)

		job := newDeployGroupJob(dg.Name, dgWf)
		if prevGroupJob != nil {
			job.Needs = []string{prevGroupJob.Name}
		}
		prevGroupJob = &job
		top.Jobs[dg.Name] = job
	}

	workflows = append(workflows, top)

	return workflows
}

func newDeploymentGroupWorkflows(group v1alpha1.DestinationGroup) (gocto.Workflow, []gocto.Workflow) {
	var subWorkflows []gocto.Workflow

	groupWf := gocto.Workflow{
		Name: group.Name,
		Jobs: make(map[string]gocto.Job),
	}

	for _, dest := range group.Destinations {
		wf := newDeploymentWorkflow(dest)
		job := newDeployJob(newDeployJobInput{
			name:           dest.FriendlyName,
			destination:    dest,
			checkoutBranch: "",
			argoCDLoginURL: "",
		})
		groupWf.Jobs[dest.FriendlyName] = job
		subWorkflows = append(subWorkflows, wf)
	}

	return groupWf, subWorkflows
}

func newDeploymentWorkflow(destination v1alpha1.Destination) gocto.Workflow {
	jobName := destination.FriendlyName
	job := newDeployJob(newDeployJobInput{
		name:           jobName,
		destination:    destination,
		checkoutBranch: "",
		argoCDLoginURL: "",
	})

	wf := gocto.Workflow{
		Name: destination.FriendlyName,
		On:   gocto.WorkflowOn{},
		Concurrency: gocto.Concurrency{
			Group:            destination.FriendlyName,
			CancelInProgress: false,
		},
		Defaults: gocto.Defaults{
			Run: gocto.DefaultsRun{
				Shell: gocto.ShellBash,
			},
		},
		Jobs: map[string]gocto.Job{
			jobName: job,
		},
	}

	return wf
}
