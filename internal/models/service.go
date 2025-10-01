package models

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
	DestinationGroups     []DestinationGroup                `json:"destinationGroups"`
	DestinationNamespace  string                            `json:"destinationNamespace"`
	SyncPolicy            *argov1alpha1.SyncPolicy          `json:"syncPolicy,omitempty,omitzero"`
}

type Source struct {
	Path    string                                `json:"path"`
	Include string                                `json:"include,omitempty,omitzero"`
	Exclude string                                `json:"exclude,omitempty,omitzero"`
	Jsonnet argov1alpha1.ApplicationSourceJsonnet `json:"jsonnet,omitempty,omitzero"`
}
