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
