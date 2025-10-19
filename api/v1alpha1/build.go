package v1alpha1

import (
	"fmt"

	"github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	"github.com/ghostsquad/alveus/internal/util"
	"github.com/goccy/go-yaml"
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
	for gIdx, group := range s.DestinationGroups {
		for dIdx, dest := range group.Destinations {
			if dest.ApplicationDestination == nil {
				dest.ApplicationDestination = &v1alpha1.ApplicationDestination{}
				s.DestinationGroups[gIdx].Destinations[dIdx] = dest
			}

			dest.Namespace = util.CoalesceStrings(
				dest.Namespace,
				group.DestinationNamespace,
				s.DestinationNamespace,
			)

			dest.ArgoCDLogin.Hostname = util.CoalesceStrings(
				dest.ArgoCDLogin.Hostname,
				group.ArgoCDLogin.Hostname,
				s.ArgoCDLogin.Hostname,
			)
		}
	}
}
