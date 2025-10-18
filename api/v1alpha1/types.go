package v1alpha1

import (
	"errors"
	"fmt"

	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	rolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
)

type Service struct {
	Name                              string                            `json:"name"`
	Source                            Source                            `json:"source"`
	IgnoreDifferences                 argov1alpha1.IgnoreDifferences    `json:"ignoreDifferences,omitempty,omitzero"`
	PrePromotionAnalysis              *rolloutsv1alpha1.RolloutAnalysis `json:"prePromotionAnalysis,omitempty,omitzero"`
	PostPromotionAnalysis             *rolloutsv1alpha1.RolloutAnalysis `json:"postPromotionAnalysis,omitempty,omitzero"`
	DestinationGroups                 []DestinationGroup                `json:"destinationGroups"`
	DestinationNamespace              string                            `json:"destinationNamespace"`
	SyncPolicy                        *argov1alpha1.SyncPolicy          `json:"syncPolicy,omitempty,omitzero"`
	ApplicationNameUniquenessStrategy ApplicationNameUniquenessStrategy `json:"applicationNameUniquenessStrategy,omitempty,omitzero"`
}

type ApplicationNameUniquenessStrategy struct {
	IncludeDestinationNamespace bool `json:"usingManyNamespaces,omitempty,omitzero"`
	IncludeGroup                bool `json:"includeGroup,omitempty,omitzero"`
}

func (s *Service) Validate() error {
	if s == nil {
		return errors.New("service is nil")
	}

	var errs []error

	if s.Name == "" {
		errs = append(errs, errors.New("service name is required"))
	}

	err := s.Source.Validate()
	if err != nil {
		errs = append(errs, fmt.Errorf("validating source: %w", err))
	}

	if len(s.DestinationGroups) == 0 {
		errs = append(errs, errors.New("at least 1 destination group is required"))
	}

	for _, destinationGroup := range s.DestinationGroups {
		err := destinationGroup.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("validating destination group: %w", err))
		}
	}

	return errors.Join(errs...)
}

type Source struct {
	Path    string                                `json:"path"`
	Include string                                `json:"include,omitempty,omitzero"`
	Exclude string                                `json:"exclude,omitempty,omitzero"`
	Jsonnet argov1alpha1.ApplicationSourceJsonnet `json:"jsonnet,omitempty,omitzero"`
}

func (s *Source) Validate() error {
	if s == nil {
		return errors.New("source is nil")
	}

	return nil
}

type DestinationGroup struct {
	Name                 string        `json:"name"`
	Destinations         []Destination `json:"destinations"`
	DestinationNamespace string        `json:"destinationNamespace,omitempty,omitzero"`
}

func (dg *DestinationGroup) Validate() error {
	if dg == nil {
		return errors.New("destinationGroup is nil")
	}

	var errs []error

	if dg.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}

	if len(dg.Destinations) == 0 {
		errs = append(errs, errors.New("destinations is empty"))
	}

	for _, d := range dg.Destinations {
		err := d.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("validating destination: %w", err))
		}
	}

	return errors.Join(errs...)
}

type Destination struct {
	*argov1alpha1.ApplicationDestination
	ArgoCDLogin ArgoCDLogin `json:"argocdLogin,omitempty,omitzero"`
}

type ArgoCDLogin struct {
	URL string `json:"url,omitempty,omitzero"`
}

func (d *Destination) Validate() error {
	if d == nil {
		return errors.New("destination is nil")
	}

	var errs []error

	if d.Namespace == "" {
		errs = append(errs, errors.New("destination namespace is required"))
	}

	if d.Name == "" && d.Server == "" {
		errs = append(errs, errors.New("one of clusterName or clusterUrl required"))
	}

	if d.Name != "" && d.Server != "" {
		errs = append(errs, errors.New("only one of clusterName or clusterUrl may be specified"))
	}

	return errors.Join(errs...)
}
