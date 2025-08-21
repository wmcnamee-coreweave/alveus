package cmd

import (
	"fmt"
	"io"
	"os"

	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	gocto "github.com/cakehappens/gocto"
	"github.com/goforj/godump"
	"github.com/spf13/cobra"

	"github.com/ghostsquad/alveus/api/v1alpha1"
	"github.com/ghostsquad/alveus/internal/integrations/argocd"
	"github.com/ghostsquad/alveus/internal/integrations/github/jobs"
	"github.com/ghostsquad/alveus/internal/integrations/github/workflows"
	"github.com/ghostsquad/alveus/internal/util"
)

func NewGenerateCommand() *cobra.Command {
	var repoURL string
	var serviceFile string

	cmd := &cobra.Command{
		Use: "generate",
		RunE: func(cmd *cobra.Command, args []string) error {
			var serviceBytes []byte
			var err error

			if serviceFile == "" || serviceFile == "-" {
				stat, _ := os.Stdin.Stat()
				if (stat.Mode() & os.ModeCharDevice) == 0 {
					serviceBytes, err = io.ReadAll(os.Stdin)
					if err != nil {
						return fmt.Errorf("reading stdin: %w", err)
					}
				} else {
					return fmt.Errorf("stdin is from a terminal")
				}
			} else {
				serviceBytes, err = os.ReadFile(serviceFile)
				if err != nil {
					return fmt.Errorf("reading from file: %s: %w", serviceFile, err)
				}
			}

			var service v1alpha1.Service
			{
				service, err = v1alpha1.NewFromYaml(serviceBytes)
				if err != nil {
					return fmt.Errorf("constructing/validating service definition: %w", err)
				}
			}

			var apps []argov1alpha1.Application
			apps, err = generateApps(repoURL, "HEAD", service)
			if err != nil {
				return fmt.Errorf("generating apps: %w", err)
			}

			var wfs []gocto.Workflow
			wfs, err = generateWorkflows(service)
			if err != nil {
				return fmt.Errorf("generating workflows: %w", err)
			}

			godump.Dump(apps)
			godump.Dump(wfs)

			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&repoURL, "repo-url", "r", "", "URL of the repository")
	err := cmd.MarkFlagRequired("repo-url")
	if err != nil {
		panic(err)
	}

	f.StringVarP(&serviceFile, "service-file", "s", "", `path to a service file. Omit or use "-" for stdin`)

	return cmd
}

func generateApps(repoURL, targetRevision string, service v1alpha1.Service) ([]argov1alpha1.Application, error) {
	var apps []argov1alpha1.Application

	for _, group := range service.DestinationGroups {
		for _, dest := range group.Destinations {
			app, err := argocd.NewApplication(argocd.Input{
				Name:           util.Join("-", service.Name, group.Name, dest.FriendlyName),
				RepoURL:        repoURL,
				TargetRevision: targetRevision,
				Sources:        service.Sources,
			}, argocd.FromServiceAPI(service))

			if err != nil {
				return nil, fmt.Errorf("constructing application: %w", err)
			}

			apps = append(apps, app)
		}
	}

	return apps, nil
}

func generateWorkflows(service v1alpha1.Service) ([]gocto.Workflow, error) {
	var wfs []gocto.Workflow

	for _, group := range service.DestinationGroups {
		var wfJobs []gocto.Job

		for _, dest := range group.Destinations {
			job := jobs.NewStage(dest.FriendlyName)
			wfJobs = append(wfJobs, job)
		}

		wf := workflows.NewDeployment(group.Name, wfJobs...)
		wfs = append(wfs, wf)
	}

	for i, wf := range wfs {
		if i+1 <= len(wfs) {
			nextWf := wfs[i+1]
			promoJob := jobs.NewPromotion(nextWf.Name, gocto.FilenameFor(nextWf))
			for _, j := range wf.Jobs {
				promoJob.Needs = append(promoJob.Needs, j.Name)
			}
			wf.Jobs[promoJob.Name] = promoJob
			wfs[i] = wf
		}
	}

	return wfs, nil
}
