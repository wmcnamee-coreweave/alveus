package v1alpha1

import (
	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	rolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
)

type Service struct {
	Name                  string                            `json:"name"`
	Sources               []Source                          `json:"sources"`
	IgnoreDifferences     argov1alpha1.IgnoreDifferences    `json:"ignoreDifferences,omitempty,omitzero"`
	PrePromotionAnalysis  *rolloutsv1alpha1.RolloutAnalysis `json:"prePromotionAnalysis,omitempty,omitzero"`
	PostPromotionAnalysis *rolloutsv1alpha1.RolloutAnalysis `json:"postPromotionAnalysis,omitempty,omitzero"`
	DestinationGroups     DestinationGroups                 `json:"destinationGroups"`
	DestinationNamespace  string                            `json:"destinationNamespace"`
	SyncPolicy            *argov1alpha1.SyncPolicy          `json:"syncPolicy,omitempty,omitzero"`
}

type Source struct {
	Path    string
	Include string
	Exclude string
	Jsonnet argov1alpha1.ApplicationSourceJsonnet
}

type DestinationGroups []DestinationGroup

type DestinationGroup struct {
	Name                 string        `json:"name"`
	Destinations         []Destination `json:"destinations"`
	DestinationNamespace string        `json:"destinationNamespace,omitempty,omitzero"`
}

type Destination struct {
	FriendlyName string `json:"friendlyName"`
	// ClusterURL specifies the URL of the target cluster's Kubernetes control plane API. This must be set if ClusterName is not set.
	ClusterURL string `json:"clusterUrl,omitempty,omitzero"`
	// ClusterName is an alternate way of specifying the target cluster by its symbolic name. This must be set if ClusterNameURL is not set.
	ClusterName string `json:"clusterName,omitempty,omitzero"`
}

func (d Destination) String() string {
	if d.ClusterURL != "" {
		return d.ClusterURL
	} else {
		return d.ClusterName
	}
}
