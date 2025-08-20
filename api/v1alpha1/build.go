package v1alpha1

import (
	"errors"
)

//func New() (*Config, error) {
//	return &Config{}, nil
//}
//
//func NewFromYaml(contents []byte) (Config, error) {
//	config := Config{}
//	err := yaml.Unmarshal(contents, &config)
//	if err != nil {
//		return config, fmt.Errorf("unmarshalling yaml: %w", err)
//	}
//
//	return config, config.Validate()
//}
//
//func (c *Config) Validate() error {
//	var errs []error
//
//	for _, svc := range c.Services {
//		errs = append(errs, svc.Validate())
//	}
//	return errors.Join(errs...)
//}

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
