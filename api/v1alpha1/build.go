package v1alpha1

import (
	"errors"
	"fmt"

	"sigs.k8s.io/yaml"
)

func NewFromYaml(contents []byte) (Service, error) {
	service := Service{}
	err := yaml.Unmarshal(contents, &service)
	if err != nil {
		return service, fmt.Errorf("unmarshalling yaml: %w", err)
	}

	return service, service.Validate()
}

func (s *Service) Validate() error {
	if s == nil {
		return errors.New("service is nil")
	}

	var errs []error

	if s.Name == "" {
		errs = append(errs, errors.New("service name is required"))
	}

	if len(s.Sources) == 0 {
		errs = append(errs, errors.New("at least 1 source is required"))
	}

	for _, source := range s.Sources {
		errs = append(errs, source.Validate())
	}

	return errors.Join(errs...)
}

func (s *Source) Validate() error {
	return nil
}
