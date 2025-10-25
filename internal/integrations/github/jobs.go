package github

import (
	"fmt"

	"github.com/cakehappens/gocto"
	"github.com/goforj/godump"

	"github.com/ghostsquad/alveus/api/v1alpha1"
	"github.com/ghostsquad/alveus/internal/util"
)

func newDeployGroupJob(name string, wf gocto.Workflow) gocto.Job {
	workflowPath := "./" + wf.GetRelativePathAndFilename()

	godump.Dump(workflowPath)

	job := gocto.Job{
		Name: name,
		Uses: workflowPath,
		Secrets: &gocto.Secrets{
			Inherit: true,
		},
	}

	return job
}

type newDeployJobInput struct {
	name                 string
	destination          v1alpha1.Destination
	checkoutCommitBranch string
	appFilePath          string
	argoLoginCommandArgs []string
	syncTimeoutSeconds   int
}

func newDeployJob(input newDeployJobInput) gocto.Job {
	const (
		EnvNameArgoCDApplicationFile = "ARGOCD_APPLICATION_FILE"
		EnvNameGitCommitMessage      = "GIT_COMMIT_MESSAGE"
		EnvNameNewTargetRevision     = "ARGOCD_APPLICATION_NEW_TARGET_REVISION"
	)

	name := input.name
	destination := input.destination

	if input.syncTimeoutSeconds == 0 {
		input.syncTimeoutSeconds = 300
	}

	destinationFriendlyName := v1alpha1.CoalesceSanitizeDestination(destination)

	steps := []gocto.Step{
		{
			Uses: "checkout@v4",
			With: map[string]any{
				"ref": input.checkoutCommitBranch,
				// otherwise, the token used is the GITHUB_TOKEN, instead of your personal token
				"persist-credentials": false,
				// otherwise, you will fail to push refs to dest repo
				"fetch-depth": 0,
			},
		},
	}

	steps = append(steps, input.destination.Github.PreDeploySteps...)

	steps = append(steps,
		gocto.Step{
			Name: "git-config",
			Run: util.SprintfDedent(`
					git config --global user.name '${{ github.actor }}'
					git config --global user.email '${{ github.actor }}@users.noreply.github.com'
				`),
		},
		gocto.Step{
			Uses: "frenck/action-setup-yq@v1",
		},
		gocto.Step{
			Name: "update-application-yaml",
			Run: util.SprintfDedent(`
					yq e '.spec.source.targetRevision = "${{ env.%s }}"' \
					'${%s}'
				`, EnvNameNewTargetRevision, EnvNameArgoCDApplicationFile),
		},
		gocto.Step{
			Name: "git-add-commit",
			Run: util.SprintfDedent(`
					git add "${%s}"
					git commit -m "${%s}"
				`, EnvNameArgoCDApplicationFile, EnvNameGitCommitMessage),
		},
		gocto.Step{
			Uses: "actions-js/push@v1.5",
			With: map[string]any{
				"github_token": "${{ secrets.GITHUB_TOKEN }}",
				"branch":       input.checkoutCommitBranch,
			},
		},
		gocto.Step{
			Name: "argocd-login",
			Run: util.SprintfDedent(`
					argocd login \
						%s \
						;
				`, util.Join(` \`+"\n\t", input.argoLoginCommandArgs...)),
		},
		gocto.Step{
			Name: "argocd-upsert",
			Run: util.SprintfDedent(`
					argocd app create \
						--grpc-web \
						--upsert \
						--file "${%s}" \
						;
				`, EnvNameArgoCDApplicationFile),
		},
		gocto.Step{
			Name: "argocd-sync",
			Run: util.SprintfDedent(`
					argocd app sync \
						--grpc-web \
						--timeout %d \
						;
				`, input.syncTimeoutSeconds),
		})

	steps = append(steps, input.destination.Github.PostDeploySteps...)

	job := gocto.Job{
		Name:   name,
		RunsOn: []string{"ubuntu-latest"},
		Defaults: gocto.Defaults{
			Run: gocto.DefaultsRun{
				Shell: gocto.ShellBash,
			},
		},
		Env: map[string]string{
			EnvNameArgoCDApplicationFile: input.appFilePath,
			EnvNameGitCommitMessage:      fmt.Sprintf("feat: 🚀 deploy to %s", destinationFriendlyName),
			EnvNameNewTargetRevision:     "123new",
		},
		Steps: steps,
	}

	return job
}
