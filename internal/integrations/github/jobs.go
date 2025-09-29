package github

import (
	"fmt"

	"github.com/cakehappens/gocto"
	"github.com/lithammer/dedent"

	"github.com/ghostsquad/alveus/api/v1alpha1"
)

func newDeployGroupJob(name string, wf gocto.Workflow) gocto.Job {
	workflowPath := wf.GetRelativePathAndFilename()

	job := gocto.Job{
		Name: name,
		Uses: workflowPath,
	}

	return job
}

func newDeployJob(name string, destination v1alpha1.Destination) gocto.Job {
	const (
		EnvNameArgoCDURL             = "ARGOCD_URL"
		EnvNameArgoCDApplicationFile = "ARGOCD_APPLICATION_FILE"
		EnvNameGitCommitMessage      = "GIT_COMMIT_MESSAGE"
		EnvNameNewTargetRevision     = "ARGOCD_APPLICATION_NEW_TARGET_REVISION"
	)

	const (
		DefaultBranch = "main"
	)

	job := gocto.Job{
		Name:   name,
		RunsOn: []string{"ubuntu-latest"},
		Defaults: gocto.Defaults{
			Run: gocto.DefaultsRun{
				Shell: gocto.ShellBash,
			},
		},
		Env: map[string]string{
			EnvNameArgoCDURL:             destination.ArgoCDLogin.URL,
			EnvNameArgoCDApplicationFile: "fake-application-file.yaml",
			EnvNameGitCommitMessage:      fmt.Sprintf("feat: 🚀 deploy to %s", destination.FriendlyName),
			EnvNameNewTargetRevision:     "123new",
		},
		Environment: gocto.Environment{
			Name: destination.FriendlyName,
		},
		Steps: []gocto.Step{
			{
				Uses: "checkout@v4",
				With: map[string]any{
					"ref": DefaultBranch,
					// otherwise, the token used is the GITHUB_TOKEN, instead of your personal token
					"persist-credentials": false,
					// otherwise, you will fail to push refs to dest repo
					"fetch-depth": 0,
				},
			},
			{
				Name: "git-config",
				Run: dedent.Dedent(`
					git config --global user.name "${{ github.actor }}"
    			git config --global user.email "${{ github.actor }}@users.noreply.github.com"
				`),
			},
			{
				Uses: "frenck/action-setup-yq@v1",
			},
			{
				Name: "update-application-yaml",
				Run: fmt.Sprintf(dedent.Dedent(`
					yq e '.spec.source.targetRevision = "${{ env.%s }}"' "${%s}"
				`), EnvNameNewTargetRevision, EnvNameArgoCDApplicationFile),
			},
			{
				Name: "git-add-commit",
				Run: fmt.Sprintf(dedent.Dedent(`
					git add "${%s}"
					git commit -m "${%s}"
				`), EnvNameArgoCDApplicationFile, EnvNameGitCommitMessage),
			},
			{
				Uses: "actions-js/push@v1.5",
				With: map[string]any{
					"github_token": "${{ secrets.GITHUB_TOKEN }}",
					"branch":       DefaultBranch,
				},
			},
			{
				Name: "argocd-login",
				Run: fmt.Sprintf(dedent.Dedent(`
					argocd login "${%s}" --grpc-web --skip-test-tls
				`), EnvNameArgoCDURL),
			},
			{
				Name: "argocd-upsert",
				Run: fmt.Sprintf(dedent.Dedent(`
					argocd app create --upsert --file "${%s}" \
						--grpc-web \
						--sync-retry-backoff-factor 2 \
						--sync-retry-backoff-max-duration 3m0s \
						--sync-retry-limit 2 \
						;
				`), EnvNameArgoCDApplicationFile),
			},
		},
	}

	return job
}
