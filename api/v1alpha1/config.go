package v1alpha1

import (
	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	rolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
)

type Config struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name                  string                            `json:"name"`
	Source                Source                            `json:"sources"`
	IgnoreDifferences     argov1alpha1.IgnoreDifferences    `json:"ignoreDifferences,omitempty"`
	PrePromotionAnalysis  *rolloutsv1alpha1.RolloutAnalysis `json:"prePromotionAnalysis,omitempty"`
	PostPromotionAnalysis *rolloutsv1alpha1.RolloutAnalysis `json:"postPromotionAnalysis,omitempty"`
	DestinationGroups     DestinationGroups                 `json:"destinations"`
	DestinationNamespace  string                            `json:"destinationNamespace"`
	SyncPolicy            *argov1alpha1.SyncPolicy          `json:"syncPolicy,omitempty"`
}

type DestinationGroups []DestinationGroup

type DestinationGroup struct {
	Name         string        `json:"name"`
	Destinations []Destination `json:"destinations"`
}

type Destination struct {
	// ClusterURL specifies the URL of the target cluster's Kubernetes control plane API. This must be set if ClusterName is not set.
	ClusterURL string `json:"clusterUrl,omitempty"`
	// ClusterName is an alternate way of specifying the target cluster by its symbolic name. This must be set if ClusterNameURL is not set.
	ClusterName string `json:"clusterName,omitempty"`
}

func (d Destination) String() string {
	if d.ClusterURL != "" {
		return d.ClusterURL
	} else {
		return d.ClusterName
	}
}

type Source struct {
	Path string `json:"path,omitempty"`
}
