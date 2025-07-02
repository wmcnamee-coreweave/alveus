// models package
// notes:
// https://github.com/argoproj/argo-cd/blob/04794332d28e44c18c89b83e9c54f0bcd69bcc43/cmd/argocd/commands/app.go#L1985
// https://github.com/marketplace/actions/manual-workflow-approval

package models

import "github.com/ghostsquad/alveus/internal/github"

// Stage is an abstraction around a GitHub workflow
// the workflow updates the associated ArgoCD Application yaml file in the /state directory
// the same changes are also pushed to ArgoCD to trigger the application to sync
type Stage struct{}

func (s *Stage) ToWorkflow() (*github.Workflow, error) {
	wf := github.Workflow{
		Name:        "",
		RunName:     "",
		On:          github.WorkflowOn{},
		Concurrency: github.Concurrency{},
		Defaults:    github.Defaults{},
		Jobs: map[string]github.Job{
			"run": {
				Name:        "",
				Permissions: github.Permissions{},
				Needs:       nil,
				If:          "",
				RunsOn:      nil,
				Environment: github.Environment{
					Name: "",
					URL:  "",
				},
				Concurrency: github.Concurrency{},
				Outputs:     nil,
				Env:         nil,
				Defaults:    github.Defaults{},
				Steps: []github.Step{
					{
						ID:               "",
						Name:             "",
						If:               "",
						Uses:             "",
						Run:              "alveus",
						WorkingDirectory: "",
						Shell:            "",
						With:             nil,
						Env:              nil,
						ContinueOnError:  false,
						TimeoutMinutes:   0,
					},
				},
				TimeoutMinutes:  0,
				ContinueOnError: false,
				Uses:            "",
				With:            nil,
				Secrets:         github.Secrets{},
			},
		},
	}

	return &wf, nil
}
