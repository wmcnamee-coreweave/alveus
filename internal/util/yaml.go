package util

import (
	"fmt"
	"reflect"

	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	"github.com/goccy/go-yaml"
)

func YamlMarshalWithOptions(val any) ([]byte, error) {
	if reflect.TypeOf(val) == reflect.TypeFor[argov1alpha1.Application]() {
		pass1, err := yaml.Marshal(val)
		if err != nil {
			return nil, fmt.Errorf("marshalling to yaml (pass 1): %w", err)
		}

		// TODO unmarshalling to a map messes up ordering
		// but there's some incompatibility between yaml and github.com/wk8/go-ordered-map/v2
		// yaml cannot unmarshal into an ordered map from that package
		pass2 := make(map[string]any)
		err = yaml.Unmarshal(pass1, &pass2)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling to map (pass 1.5): %w", err)
		}

		CleanKubernetesManifest(pass2)
		val = pass2
	}

	result, err := yaml.Marshal(val)
	if err != nil {
		return nil, fmt.Errorf("marshalling yaml (pass2): %w", err)
	}
	return result, nil
}

func CleanKubernetesManifest(val map[string]any) {
	if metadata, ok := val["metadata"]; ok {
		if metadataMap, ok := metadata.(map[string]any); ok {
			delete(metadataMap, "creationTimestamp")
		}

		delete(val, "status")
	}
}
