package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	"github.com/cakehappens/gocto"
	"github.com/go-git/go-billy/v6"
	"github.com/go-git/go-billy/v6/osfs"
	billyutil "github.com/go-git/go-billy/v6/util"
	"github.com/spf13/cobra"

	"github.com/ghostsquad/alveus/api/v1alpha1"
	"github.com/ghostsquad/alveus/internal/constants"
	"github.com/ghostsquad/alveus/internal/integrations/argocd"
	"github.com/ghostsquad/alveus/internal/integrations/github"
	"github.com/ghostsquad/alveus/internal/util"
)

func NewGenerateCommand() *cobra.Command {
	var repoURL string
	var serviceFile string
	var applicationOutputPath string
	var workflowOutputPath string

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

			wfs = github.NewWorkflows(service)

			{
				fs := osfs.New(".")
				if err := writeApps(fs, applicationOutputPath, apps); err != nil {
					return fmt.Errorf("writing apps: %w", err)
				}
				if err := writeWorkflows(fs, workflowOutputPath, wfs); err != nil {
					return fmt.Errorf("writing workflows: %w", err)
				}
			}

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
	f.StringVar(&applicationOutputPath, "application-output-path", "./.alveus/applications", "path to where to write ArgoCD application resources")

	f.StringVar(&workflowOutputPath, "workflow-output-path", gocto.DefaultPathToWorkflows, "path to where to write Github workflow files")

	return cmd
}

func generateApps(repoURL, targetRevision string, service v1alpha1.Service) ([]argov1alpha1.Application, error) {
	var apps []argov1alpha1.Application

	for _, group := range service.DestinationGroups {
		for _, dest := range group.Destinations {
			app, err := argocd.NewApplication(argocd.Input{
				Name:                   util.Join("-", service.Name, group.Name, dest.FriendlyName),
				RepoURL:                repoURL,
				TargetRevision:         targetRevision,
				DestinationNamespace:   service.DestinationNamespace,
				DestinationClusterName: dest.ClusterName,
				DestinationClusterURL:  dest.ClusterURL,
			}, argocd.FromServiceAPI(service))

			if err != nil {
				return nil, fmt.Errorf("constructing application: %w", err)
			}

			apps = append(apps, app)
		}
	}

	return apps, nil
}

func writeApps(fs billy.Filesystem, basepath string, apps []argov1alpha1.Application) error {
	if err := fs.MkdirAll(basepath, os.ModePerm); err != nil {
		return fmt.Errorf("creating directory: %q: %w", basepath, err)
	}

	if err := billyutil.RemoveAll(fs, basepath); err != nil {
		return fmt.Errorf("cleaning directory: %q: %w", basepath, err)
	}

	for _, app := range apps {
		filename := argocd.FilenameFor(app)
		fullFilename := filepath.Join(basepath, filename)
		fileBytes, err := util.YamlMarshalWithOptions(app)
		if err != nil {
			return fmt.Errorf("marshalling application to yaml: %w", err)
		}

		if err := billyutil.WriteFile(fs, fullFilename, fileBytes, os.ModePerm); err != nil {
			return fmt.Errorf("writing application to file: %q: %w", fullFilename, err)
		}
	}

	return nil
}

func writeWorkflows(fs billy.Filesystem, basepath string, wfs []gocto.Workflow) error {
	if err := fs.MkdirAll(basepath, os.ModePerm); err != nil {
		return fmt.Errorf("creating directory: %q: %w", basepath, err)
	}

	files, err := fs.ReadDir(basepath)
	if err != nil {
		return fmt.Errorf("reading directory: %q: %w", basepath, err)
	}

	expectedPrefix := constants.Alveus + "-"

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), expectedPrefix) {
			if err := fs.Remove(filepath.Join(basepath, file.Name())); err != nil {
				return fmt.Errorf("removing file: %q: %w", file.Name(), err)
			}
		}
	}

	for _, wf := range wfs {
		filename := expectedPrefix + wf.GetFilename()
		fullFilename := filepath.Join(basepath, filename)
		fileBytes, err := util.YamlMarshalWithOptions(wf)
		if err != nil {
			return fmt.Errorf("marshalling workflow to yaml: %w", err)
		}

		if err := billyutil.WriteFile(fs, fullFilename, fileBytes, os.ModePerm); err != nil {
			return fmt.Errorf("writing workflow to file: %q: %w", fullFilename, err)
		}
	}

	return nil
}
