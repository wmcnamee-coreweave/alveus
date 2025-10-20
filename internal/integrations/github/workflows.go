package github

import (
	"github.com/cakehappens/gocto"

	"github.com/ghostsquad/alveus/api/v1alpha1"
)

func NewWorkflows(service v1alpha1.Service) []gocto.Workflow {
	var workflows []gocto.Workflow

	top := gocto.Workflow{
		Name: service.Name,
		On:   service.Github.On,
		Jobs: make(map[string]gocto.Job),
	}

	var prevGroupJob *gocto.Job
	for _, dg := range service.DestinationGroups {
		dgWf, subWfs := newDeploymentGroupWorkflows(newDeploymentGroupWorkflowInput{
			namePrefix:           service.Name,
			group:                dg,
			checkoutCommitBranch: service.ArgoCD.Source.CommitBranch,
		})
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

type newDeploymentGroupWorkflowInput struct {
	namePrefix           string
	group                v1alpha1.DestinationGroup
	checkoutCommitBranch string
}

func newDeploymentGroupWorkflows(input newDeploymentGroupWorkflowInput) (gocto.Workflow, []gocto.Workflow) {
	var subWorkflows []gocto.Workflow

	groupWf := gocto.Workflow{
		Name: input.namePrefix + "-" + input.group.Name,
		On: gocto.WorkflowOn{
			Dispatch: &gocto.OnDispatch{},
		},
		Jobs: make(map[string]gocto.Job),
	}

	for _, dest := range input.group.Destinations {
		wf := newDeploymentWorkflow(newDeploymentWorkflowInput{
			namePrefix:           input.namePrefix,
			checkoutCommitBranch: input.checkoutCommitBranch,
			destination:          dest,
		})
		destinationFriendlyName := v1alpha1.CoalesceSanitizeDestination(dest)
		groupWf.Jobs[destinationFriendlyName] = newDeployGroupJob(destinationFriendlyName, wf)
		subWorkflows = append(subWorkflows, wf)
	}

	return groupWf, subWorkflows
}

type newDeploymentWorkflowInput struct {
	namePrefix           string
	checkoutCommitBranch string
	destination          v1alpha1.Destination
}

func newDeploymentWorkflow(input newDeploymentWorkflowInput) gocto.Workflow {
	destinationFriendlyName := v1alpha1.CoalesceSanitizeDestination(input.destination)

	jobName := destinationFriendlyName
	job := newDeployJob(newDeployJobInput{
		name:                 jobName,
		destination:          input.destination,
		checkoutCommitBranch: input.checkoutCommitBranch,
		argocdHostname:       input.destination.ArgoCD.Hostname,
	})

	wf := gocto.Workflow{
		Name: input.namePrefix + "-" + destinationFriendlyName,
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
