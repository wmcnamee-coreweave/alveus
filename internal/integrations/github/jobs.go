package github

import (
	"github.com/cakehappens/gocto"
	"github.com/lithammer/dedent"

	"github.com/ghostsquad/alveus/api/v1alpha1"
)

func newDeployGroupJob(name string, wf gocto.Workflow) gocto.Job {
	workflowPath := wf.GetRelativePathAndFilename()

	job := gocto.Job{
		Name:        name,
		Permissions: gocto.Permissions{},
		Uses:        workflowPath,
		With:        map[string]any{},
	}

	return job
}

func newDeployJob(name string, destination v1alpha1.Destination) gocto.Job {
	job := gocto.Job{
		Name:   name,
		RunsOn: []string{"ubuntu-latest"},
		Defaults: gocto.Defaults{
			Run: gocto.DefaultsRun{
				Shell: gocto.ShellBash,
			},
		},
		Environment: gocto.Environment{
			Name: destination.FriendlyName,
		},
		Steps: []gocto.Step{
			{
				Uses: "checkout@v4",
			},
			{
				Run: dedent.Dedent(`
					argocd login 
				`),
			},
			{
				Run: dedent.Dedent(`
					argocd app create --upsert
				`),
			},
		},
	}

	return job
}
