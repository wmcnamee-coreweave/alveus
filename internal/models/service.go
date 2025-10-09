package models

import (
	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	rolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"

	"github.com/ghostsquad/alveus/api/v1alpha1"
)

type Service struct {
	Name                  string
	Source                Source
	IgnoreDifferences     argov1alpha1.IgnoreDifferences
	PrePromotionAnalysis  *rolloutsv1alpha1.RolloutAnalysis
	PostPromotionAnalysis *rolloutsv1alpha1.RolloutAnalysis
	DestinationGroups     []DestinationGroup
	DestinationNamespace  string
	SyncPolicy            *argov1alpha1.SyncPolicy
}

func NewDestinationGroupsFromV1Alpha1(d []v1alpha1.DestinationGroup) []DestinationGroup {
	destinationGroups := make([]DestinationGroup, len(d))
	for i := range d {
		destinationGroups[i] = NewDestinationGroupFromV1Alpha1(d[i])
	}

	return destinationGroups
}

func NewDestinationGroupFromV1Alpha1(d v1alpha1.DestinationGroup) DestinationGroup {
	destinations := make([]Destination, len(d.Destinations))
	for i := range d.Destinations {
		destinations[i] = NewDestinationFromV1Alpha1(d.Destinations[i])
	}

	return DestinationGroup{
		Name:                 d.Name,
		Destinations:         destinations,
		DestinationNamespace: d.DestinationNamespace,
	}
}

func NewDestinationFromV1Alpha1(d v1alpha1.Destination) Destination {
	return Destination{
		Name:                  "",
		FriendlyName:          d.FriendlyName,
		PrePromotionAnalysis:  nil,
		PostPromotionAnalysis: nil,
		ArgoCD:                nil,
		GitHub:                nil,
	}
}

func NewServiceFromV1Alpha1(service v1alpha1.Service) Service {
	return Service{
		Name:                  service.Name,
		Source:                NewSourceFromV1Alpha1(service.Source),
		IgnoreDifferences:     service.IgnoreDifferences,
		PrePromotionAnalysis:  service.PrePromotionAnalysis,
		PostPromotionAnalysis: service.PostPromotionAnalysis,
		DestinationGroups:     NewDestinationGroupsFromV1Alpha1(service.DestinationGroups),
		DestinationNamespace:  "",
		SyncPolicy:            nil,
	}
}

type Source struct {
	Path    string                                `json:"path"`
	Include string                                `json:"include,omitempty,omitzero"`
	Exclude string                                `json:"exclude,omitempty,omitzero"`
	Jsonnet argov1alpha1.ApplicationSourceJsonnet `json:"jsonnet,omitempty,omitzero"`
}

func NewSourceFromV1Alpha1(source v1alpha1.Source) Source {
	return Source{
		Path:    source.Path,
		Include: source.Include,
		Exclude: source.Exclude,
		Jsonnet: source.Jsonnet,
	}
}
