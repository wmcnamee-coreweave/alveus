package github

import (
	"fmt"
	"github.com/cakehappens/gocto"
	"github.com/goforj/godump"

	"github.com/wmcnamee-coreweave/alveus/api/v1alpha1"
	"github.com/wmcnamee-coreweave/alveus/internal/constants"
	"github.com/wmcnamee-coreweave/alveus/internal/integrations/argocd"
	"github.com/wmcnamee-coreweave/alveus/internal/util"
)

func SetWorkflowFilenameWithAlveusPrefix(w gocto.Workflow) gocto.Workflow {
	filename := gocto.FilenameFor(w)
	filename = constants.Alveus + "-" + filename
	(&w).SetFilename(filename)

	return w
}

func NewWorkflows(service v1alpha1.Service, apps argocd.ApplicationRepository) []gocto.Workflow {
	var workflows []gocto.Workflow

	top := gocto.Workflow{
		Name: service.Name,
		On:   service.Github.On,
		Jobs: make(map[string]gocto.Job),
	}

	top = SetWorkflowFilenameWithAlveusPrefix(top)

	var prevGroupJob *gocto.Job
	for _, dg := range service.DestinationGroups {
		dgWf, subWfs := newDeploymentGroupWorkflows(newDeploymentGroupWorkflowInput{
			namePrefix:           service.Name,
			group:                dg,
			checkoutCommitBranch: service.ArgoCD.Source.CommitBranch,
			apps:                 apps,
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
	apps                 argocd.ApplicationRepository
}

func newDeploymentGroupWorkflows(input newDeploymentGroupWorkflowInput) (gocto.Workflow, []gocto.Workflow) {
	var subWorkflows []gocto.Workflow

	groupWf := gocto.Workflow{
		Name: input.namePrefix + "-" + input.group.Name,
		On: gocto.WorkflowOn{
			Dispatch: &gocto.OnDispatch{},
			Call:     &gocto.OnCall{},
		},
		Jobs: make(map[string]gocto.Job),
	}
	groupWf = SetWorkflowFilenameWithAlveusPrefix(groupWf)

	for _, dest := range input.group.Destinations {
		wf := newDeploymentWorkflow(newDeploymentWorkflowInput{
			namePrefix:           input.namePrefix,
			checkoutCommitBranch: input.checkoutCommitBranch,
			destination:          dest,
			apps:                 input.apps,
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
	apps                 argocd.ApplicationRepository
}

func newDeploymentWorkflow(input newDeploymentWorkflowInput) gocto.Workflow {
	destinationFriendlyName := v1alpha1.CoalesceSanitizeDestination(input.destination)

	jobName := destinationFriendlyName

	appFilePath, _, ok := input.apps.GetByDestination(input.destination)
	if !ok {
		godump.Dump(input.apps)
		panic(fmt.Errorf("no app found for destination %+v", input.destination))
	}

	job := newDeployJob(newDeployJobInput{
		name:                 jobName,
		destination:          input.destination,
		checkoutCommitBranch: input.checkoutCommitBranch,
		argoCDSpec:           input.destination.ArgoCD,
		appFilePath:          appFilePath,
	})

	jobs := util.MergeMapsShallow(
		input.destination.Github.ExtraDeployJobs,
		map[string]gocto.Job{
			jobName: job,
		},
	)

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
		Jobs: jobs,
	}

	wf = SetWorkflowFilenameWithAlveusPrefix(wf)

	return wf
}
