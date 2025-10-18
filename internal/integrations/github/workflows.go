package github

import (
	"github.com/cakehappens/gocto"

	"github.com/ghostsquad/alveus/api/v1alpha1"
	"github.com/ghostsquad/alveus/internal/integrations/argocd"
)

func NewWorkflows(service v1alpha1.Service) []gocto.Workflow {
	var workflows []gocto.Workflow

	top := gocto.Workflow{
		Name: service.Name,
		On: gocto.WorkflowOn{
			Dispatch: &gocto.OnDispatch{},
			Call:     &gocto.OnCall{},
		},
		Jobs: make(map[string]gocto.Job),
	}

	var prevGroupJob *gocto.Job
	for _, dg := range service.DestinationGroups {
		dgWf, subWfs := newDeploymentGroupWorkflows(service.Name, dg)
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

func newDeploymentGroupWorkflows(namePrefix string, group v1alpha1.DestinationGroup) (gocto.Workflow, []gocto.Workflow) {
	var subWorkflows []gocto.Workflow

	groupWf := gocto.Workflow{
		Name: namePrefix + "-" + group.Name,
		On: gocto.WorkflowOn{
			Dispatch: &gocto.OnDispatch{},
		},
		Jobs: make(map[string]gocto.Job),
	}

	for _, dest := range group.Destinations {
		wf := newDeploymentWorkflow(namePrefix, dest)
		destinationFriendlyName := argocd.CoalesceSanitizeDestination(*dest.ApplicationDestination)
		groupWf.Jobs[destinationFriendlyName] = newDeployGroupJob(destinationFriendlyName, wf)
		subWorkflows = append(subWorkflows, wf)
	}

	return groupWf, subWorkflows
}

func newDeploymentWorkflow(namePrefix string, destination v1alpha1.Destination) gocto.Workflow {
	destinationFriendlyName := argocd.CoalesceSanitizeDestination(*destination.ApplicationDestination)

	jobName := destinationFriendlyName
	job := newDeployJob(newDeployJobInput{
		name:           jobName,
		destination:    destination,
		checkoutBranch: "",
		argoCDLoginURL: "",
	})

	wf := gocto.Workflow{
		Name: namePrefix + "-" + destinationFriendlyName,
		On: gocto.WorkflowOn{
			Dispatch: &gocto.OnDispatch{},
			Call:     &gocto.OnCall{},
		},
		Concurrency: gocto.Concurrency{
			Group:            destinationFriendlyName,
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
