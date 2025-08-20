package argocd

import (
	"errors"
	"fmt"

	argoapisapplication "github.com/argoproj/argo-cd/v3/pkg/apis/application"
	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/ghostsquad/alveus/api/v1alpha1"
)

type Options struct {
	Namespace         string
	Labels            map[string]string
	Annotations       map[string]string
	IgnoreDifferences argov1alpha1.IgnoreDifferences
	SyncPolicy        *argov1alpha1.SyncPolicy
	Project           string
}

type Option func(*Options)

type Input struct {
	Name           string
	RepoURL        string
	TargetRevision string
	Sources        []v1alpha1.Source
}

func (in *Input) Validate() error {
	var errs []error

	if in.Name == "" {
		errs = append(errs, fmt.Errorf("name is required"))
	}

	if in.RepoURL == "" {
		errs = append(errs, fmt.Errorf("repoURL is required"))
	}

	if in.TargetRevision == "" {
		errs = append(errs, fmt.Errorf("targetRevision is required"))
	}

	if len(in.Sources) == 0 {
		errs = append(errs, fmt.Errorf("at least one source is required"))
	}

	for _, source := range in.Sources {
		if source.Path == "" {
			errs = append(errs, fmt.Errorf("sources[*].path is required"))
		}
	}

	return errors.Join(errs...)
}

func FromServiceAPI(service v1alpha1.Service) Option {
	return func(o *Options) {
		o.SyncPolicy = service.SyncPolicy
		o.IgnoreDifferences = service.IgnoreDifferences
		o.Namespace = service.DestinationNamespace
	}
}

func NewApplication(input Input, options ...Option) (argov1alpha1.Application, error) {
	opts := &Options{
		Namespace:   "argocd",
		Labels:      map[string]string{},
		Annotations: map[string]string{},
		Project:     "default",
	}
	for _, o := range options {
		o(opts)
	}

	labels := opts.Labels
	if len(labels) == 0 {
		labels = nil
	}

	annotations := opts.Annotations
	if len(annotations) == 0 {
		annotations = nil
	}

	err := input.Validate()
	if err != nil {
		return nil, err
	}

	sources := make([]argov1alpha1.ApplicationSource, len(input.Sources))
	for i, source := range input.Sources {
		if source.Include == "" {
			source.Include = "{*.yml,*.yaml}"
		}

		sources[i] = argov1alpha1.ApplicationSource{
			RepoURL:        input.RepoURL,
			Path:           source.Path,
			TargetRevision: input.TargetRevision,
			Directory: &argov1alpha1.ApplicationSourceDirectory{
				Recurse: true,
				Jsonnet: source.Jsonnet,
				Exclude: source.Exclude,
				Include: source.Include,
			},
		}
	}

	app := argov1alpha1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       argoapisapplication.ApplicationKind,
			APIVersion: argov1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        input.Name,
			Namespace:   opts.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: argov1alpha1.ApplicationSpec{
			Destination:       argov1alpha1.ApplicationDestination{},
			Project:           opts.Project,
			SyncPolicy:        opts.SyncPolicy,
			IgnoreDifferences: opts.IgnoreDifferences,
			Sources:           sources,
		},
	}

	return app, nil
}
