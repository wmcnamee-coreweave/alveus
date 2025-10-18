package github

import (
	"fmt"

	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	"github.com/cakehappens/gocto"

	"github.com/ghostsquad/alveus/api/v1alpha1"
	"github.com/ghostsquad/alveus/internal/integrations/argocd"
	"github.com/ghostsquad/alveus/internal/util"
)

func newDeployGroupJob(name string, wf gocto.Workflow) gocto.Job {
	workflowPath := wf.GetRelativePathAndFilename()

	job := gocto.Job{
		Name: name,
		Uses: workflowPath,
	}

	return job
}

type newDeployJobInput struct {
	name               string
	destination        v1alpha1.Destination
	checkoutBranch     string
	argoCDLoginURL     string
	argoCDApplication  argov1alpha1.Application
	syncTimeoutSeconds int
}

func newDeployJob(input newDeployJobInput) gocto.Job {
	const (
		EnvNameArgoCDURL             = "ARGOCD_URL"
		EnvNameArgoCDApplicationFile = "ARGOCD_APPLICATION_FILE"
		EnvNameGitCommitMessage      = "GIT_COMMIT_MESSAGE"
		EnvNameNewTargetRevision     = "ARGOCD_APPLICATION_NEW_TARGET_REVISION"
		EnvNameArgoCDAuthToken       = "ARGOCD_AUTH_TOKEN"
	)

	name := input.name
	destination := input.destination

	if input.syncTimeoutSeconds == 0 {
		input.syncTimeoutSeconds = 300
	}

	destinationFriendlyName := argocd.CoalesceSanitizeDestination(*destination.ApplicationDestination)

	job := gocto.Job{
		Name:   name,
		RunsOn: []string{"ubuntu-latest"},
		Defaults: gocto.Defaults{
			Run: gocto.DefaultsRun{
				Shell: gocto.ShellBash,
			},
		},
		Env: map[string]string{
			EnvNameArgoCDURL:             input.argoCDLoginURL,
			EnvNameArgoCDApplicationFile: "fake-application-file.yaml",
			EnvNameGitCommitMessage:      fmt.Sprintf("feat: 🚀 deploy to %s", destinationFriendlyName),
			EnvNameNewTargetRevision:     "123new",
			EnvNameArgoCDAuthToken:       "fake-auth-token",
		},
		Steps: []gocto.Step{
			{
				Uses: "checkout@v4",
				With: map[string]any{
					"ref": input.checkoutBranch,
					// otherwise, the token used is the GITHUB_TOKEN, instead of your personal token
					"persist-credentials": false,
					// otherwise, you will fail to push refs to dest repo
					"fetch-depth": 0,
				},
			},
			{
				Name: "git-config",
				Run: util.SprintfDedent(`
					git config --global user.name '${{ github.actor }}'
					git config --global user.email '${{ github.actor }}@users.noreply.github.com'
				`),
			},
			{
				Uses: "frenck/action-setup-yq@v1",
			},
			{
				Name: "update-application-yaml",
				Run: util.SprintfDedent(`
					yq e '.spec.source.targetRevision = "${{ env.%s }}"' \
					'${%s}'
				`, EnvNameNewTargetRevision, EnvNameArgoCDApplicationFile),
			},
			{
				Name: "git-add-commit",
				Run: util.SprintfDedent(`
					git add "${%s}"
					git commit -m "${%s}"
				`, EnvNameArgoCDApplicationFile, EnvNameGitCommitMessage),
			},
			{
				Uses: "actions-js/push@v1.5",
				With: map[string]any{
					"github_token": "${{ secrets.GITHUB_TOKEN }}",
					"branch":       input.checkoutBranch,
				},
			},
			{
				Name: "argocd-upsert",
				Run: util.SprintfDedent(`
					argocd app create \
						--grpc-web \
						--upsert \
						--file "${%s}" \
						;
				`, EnvNameArgoCDApplicationFile),
			},
			{
				Name: "argocd-sync",
				Run: util.SprintfDedent(`
					argocd app sync \
						--grpc-web \
						--timeout %d \
						;
				`, input.syncTimeoutSeconds),
			},
		},
	}

	return job
}
