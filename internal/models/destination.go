package models

import (
	"errors"
	"fmt"

	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	rolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"

	"github.com/ghostsquad/alveus/api/v1alpha1"
	"github.com/ghostsquad/alveus/internal/integrations/argocd"
	"github.com/ghostsquad/alveus/internal/util"
)

type NewDestinationOptions struct {
	FriendlyName string
	ArgoCD       *ArgoCD
	GitHub       *GitHub
}

type NewDestinationOption func(*NewDestinationOptions) error

type ApplicationNameInput struct {
	ServiceName string
	Namespace   string
	ClusterName string
}

func NewApplicationName(in ApplicationNameInput) (string, error) {
	var errs []error
	if in.ServiceName == "" {
		errs = append(errs, errors.New("serviceName missing"))
	}

	if in.ClusterName == "" {
		errs = append(errs, errors.New("clusterName missing"))
	}

	if in.Namespace == "" {
		errs = append(errs, errors.New("namespace missing"))
	}

	err := errors.Join(errs...)
	if err != nil {
		return "", err
	}

	name := util.Join("-", in.ServiceName, in.ClusterName, in.Namespace)
	return util.SanitizeNameForKubernetes(name)
}

type SourceInfo struct {
	ServiceName    string
	RepoURL        string
	TargetRevision string
	Source         v1alpha1.Source
}

func (s *SourceInfo) Validate() error {
	var errs []error

	if s.ServiceName == "" {
		errs = append(errs, errors.New("serviceName required"))
	}

	if s.RepoURL == "" {
		errs = append(errs, errors.New("repoURL required"))
	}

	if s.TargetRevision == "" {
		errs = append(errs, errors.New("targetRevision required"))
	}

	err := errors.Join(errs...)
	if err != nil {
		return fmt.Errorf("validating sourceInfo: %w", err)
	}

	return nil
}

type DestinationInfo struct {
	ArgoCDURL   string
	Namespace   string
	ClusterName string
	ClusterURL  string
}

func (d *DestinationInfo) Validate() error {
	var errs []error

	if d.ClusterName == "" {
		errs = append(errs, errors.New("clusterName required"))
	}

	if d.Namespace == "" {
		errs = append(errs, errors.New("namespace required"))
	}

	if d.ArgoCDURL == "" {
		errs = append(errs, errors.New("argocdURL required"))
	}

	if d.ClusterName == "" && d.ClusterURL == "" {
		errs = append(errs, errors.New("clusterName or clusterURL required"))
	}

	if d.ClusterName != "" && d.ClusterURL != "" {
		errs = append(errs, errors.New("clusterName and clusterURL are mutually exclusive"))
	}

	err := errors.Join(errs...)
	if err != nil {
		return fmt.Errorf("validating destinationInfo: %w", err)
	}

	return nil
}

func FromInfos(srcInfo SourceInfo, destInfo DestinationInfo) NewDestinationOption {
	return func(opts *NewDestinationOptions) error {
		errs := []error{
			srcInfo.Validate(),
			destInfo.Validate(),
		}

		err := errors.Join(errs...)
		if err != nil {
			return err
		}

		opts.ArgoCD.LoginURL = destInfo.ArgoCDURL

		appName, err := NewApplicationName(ApplicationNameInput{
			ServiceName: srcInfo.ServiceName,
			ClusterName: destInfo.ClusterName,
			Namespace:   destInfo.Namespace,
		})
		if err != nil {
			return fmt.Errorf("constructing application name: %w", err)
		}

		app, err := argocd.NewApplication(
			argocd.Input{
				Name:                   appName,
				RepoURL:                srcInfo.RepoURL,
				TargetRevision:         srcInfo.TargetRevision,
				DestinationNamespace:   destInfo.Namespace,
				DestinationClusterName: destInfo.ClusterName,
				DestinationClusterURL:  destInfo.ClusterURL,
			},
			argocd.WithSource(srcInfo.Source),
		)
		opts.ArgoCD.Application = &app
		return err
	}
}

func NewDestination(name string, options ...NewDestinationOption) (Destination, error) {
	opts := &NewDestinationOptions{
		FriendlyName: name,
	}

	var errs []error
	for _, opt := range options {
		errs = append(errs, opt(opts))
	}

	err := errors.Join(errs...)
	if err != nil {
		return Destination{}, err
	}

	dest := Destination{
		Name:         name,
		FriendlyName: opts.FriendlyName,
		ArgoCD:       opts.ArgoCD,
		GitHub:       opts.GitHub,
	}

	return dest, nil
}

type Destination struct {
	Name                  string
	FriendlyName          string
	PrePromotionAnalysis  *rolloutsv1alpha1.RolloutAnalysis
	PostPromotionAnalysis *rolloutsv1alpha1.RolloutAnalysis
	ArgoCD                *ArgoCD
	GitHub                *GitHub
}

type ArgoCD struct {
	LoginURL    string
	Application *argov1alpha1.Application
}

type GitHub struct {
	Checkout  Checkout
	GitConfig GitConfig
}

type Checkout struct {
	Branch string
}

type GitConfig struct {
	UserName  string
	UserEmail string
}

type DestinationGroup struct {
	Name                 string
	Destinations         []Destination
	DestinationNamespace string
}
