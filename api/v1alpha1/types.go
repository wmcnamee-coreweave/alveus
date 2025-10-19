package v1alpha1

import (
	"errors"
	"fmt"
	"strings"

	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	rolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
)

type Service struct {
	Name                              string                            `json:"name"`
	Source                            Source                            `json:"source"`
	IgnoreDifferences                 argov1alpha1.IgnoreDifferences    `json:"ignoreDifferences,omitempty,omitzero"`
	PrePromotionAnalysis              *rolloutsv1alpha1.RolloutAnalysis `json:"prePromotionAnalysis,omitempty,omitzero"`
	PostPromotionAnalysis             *rolloutsv1alpha1.RolloutAnalysis `json:"postPromotionAnalysis,omitempty,omitzero"`
	DestinationGroups                 DestinationGroups                 `json:"destinationGroups"`
	DestinationNamespace              string                            `json:"destinationNamespace"`
	SyncPolicy                        *argov1alpha1.SyncPolicy          `json:"syncPolicy,omitempty,omitzero"`
	ApplicationNameUniquenessStrategy ApplicationNameUniquenessStrategy `json:"applicationNameUniquenessStrategy,omitempty,omitzero"`

	sourceValidatorFunc            func(source Source) error
	destinationGroupsValidatorFunc func(groups DestinationGroups) error
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

	if s.sourceValidatorFunc == nil {
		s.sourceValidatorFunc = func(source Source) error {
			return source.Validate()
		}
	}

	err := s.sourceValidatorFunc(s.Source)
	if err != nil {
		errs = append(errs, fmt.Errorf("validating source: %w", err))
	}

	if s.destinationGroupsValidatorFunc == nil {
		s.destinationGroupsValidatorFunc = func(groups DestinationGroups) error {
			return groups.Validate()
		}
	}

	errs = append(errs, s.destinationGroupsValidatorFunc(s.DestinationGroups))

	return errors.Join(errs...)
}

type DestinationGroups []DestinationGroup

func (dg DestinationGroups) Validate() error {
	var errs []error

	if len(dg) == 0 {
		errs = append(errs, errors.New("at least 1 destination group is required"))
	}

	groupsFound := make(map[string]struct{})

	for _, destinationGroup := range dg {
		err := destinationGroup.Validate()
		groupName := destinationGroup.Name
		if groupName == "" {
			groupName = "<empty>"
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("validating destination group: %s: %w", groupName, err))
			continue
		}

		if _, ok := groupsFound[destinationGroup.Name]; ok {
			errs = append(errs, fmt.Errorf("duplicate destination group name: %s", destinationGroup.Name))
		} else {
			groupsFound[destinationGroup.Name] = struct{}{}
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

	destinationsValidatorFunc func(destinations Destinations) error
}

func (dg *DestinationGroup) Validate() error {
	if dg == nil {
		return errors.New("destinationGroup is nil")
	}

	var errs []error

	if dg.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}

	if dg.destinationsValidatorFunc == nil {
		dg.destinationsValidatorFunc = func(ds Destinations) error {
			return ds.Validate()
		}
	}

	dsValidationErr := dg.destinationsValidatorFunc(dg.Destinations)
	if dsValidationErr != nil {

		errs = append(errs, fmt.Errorf("validating destinations: %w", dsValidationErr))
	}

	return errors.Join(errs...)
}

type Destinations []Destination

func (ds Destinations) Validate() error {
	var errs []error

	if len(ds) == 0 {
		errs = append(errs, errors.New("destinations is empty"))
	}

	destinationsFound := make(map[string]struct{})

	for _, d := range ds {
		if d.ApplicationDestination == nil {
			errs = append(errs, errors.New("destination is nil"))
			continue
		}

		err := d.Validate()
		if err != nil {
			errs = append(errs, fmt.Errorf("validating destination: %w", err))
		}

		destFinalName := CoalesceSanitizeDestination(*d.ApplicationDestination) + "/namespace/" + d.Namespace
		if _, ok := destinationsFound[destFinalName]; ok {
			errs = append(errs, fmt.Errorf("duplicate destination name: %s", destFinalName))
		} else {
			destinationsFound[destFinalName] = struct{}{}
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

func CoalesceSanitizeDestination(destination argov1alpha1.ApplicationDestination) string {
	if destination.Name != "" {
		return strings.ToLower(destination.Name)
	}

	name := strings.ToLower(destination.Server)
	name = strings.TrimPrefix(name, "https://")
	name = strings.TrimPrefix(name, "http://")
	name = strings.ReplaceAll(name, ".", "-")

	return name
}
