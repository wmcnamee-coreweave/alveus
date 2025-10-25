package v1alpha1

import (
	"fmt"

	"github.com/cakehappens/gocto"
	"github.com/goccy/go-yaml"

	"github.com/ghostsquad/alveus/internal/util"
)

func NewFromYaml(contents []byte) (Service, error) {
	service := &Service{}
	err := yaml.Unmarshal(contents, service)
	if err != nil {
		return *service, fmt.Errorf("unmarshalling yaml: %w", err)
	}

	service.Inflate()

	return *service, service.Validate()
}

func (s *Service) Inflate() {
	if s.Github.On.Dispatch == nil {
		s.Github.On.Dispatch = &gocto.OnDispatch{}
	}

	for gIdx, group := range s.DestinationGroups {
		for dIdx, dest := range group.Destinations {
			dest.Namespace = util.CoalesceStrings(
				dest.Namespace,
				group.DestinationNamespace,
				s.DestinationNamespace,
			)

			dest.ArgoCD.LoginCommandArgs = util.CoalesceSlices(
				dest.ArgoCD.LoginCommandArgs,
				group.ArgoCD.LoginCommandArgs,
				s.ArgoCD.LoginCommandArgs,
			)

			dest.Github.Secrets = util.CoalescePointers(
				dest.Github.Secrets,
				group.Github.Secrets,
				s.Github.Secrets,
				&gocto.Secrets{
					Inherit: true,
				},
			)

			dest.Github.PreDeploySteps = util.CoalesceSlices(
				dest.Github.PreDeploySteps,
				group.Github.PreDeploySteps,
				s.Github.PreDeploySteps,
			)

			dest.Github.PostDeploySteps = util.CoalesceSlices(
				dest.Github.PostDeploySteps,
				group.Github.PostDeploySteps,
				s.Github.PostDeploySteps,
			)

			dest.Github.ExtraDeployJobs = util.CoalesceMaps(
				dest.Github.ExtraDeployJobs,
				group.Github.ExtraDeployJobs,
				s.Github.ExtraDeployJobs,
			)

			s.DestinationGroups[gIdx].Destinations[dIdx] = dest
		}
	}
}
