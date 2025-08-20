package main

import (
	"context"
	"fmt"
	"os"
	"syscall"
	
	"github.com/oklog/run"
	"sigs.k8s.io/yaml"

	"github.com/ghostsquad/alveus/api/v1alpha1"
)

var version string

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var g run.Group
	var err error

	g.Add(func() error {
		service := v1alpha1.Service{
			Name: "example-service",
			Sources: []v1alpha1.Source{
				{
					Path: "./examples/example-service/manifests",
				},
			},
			DestinationGroups: []v1alpha1.DestinationGroup{
				{
					Name: "staging",
					Destinations: []v1alpha1.Destination{
						{
							FriendlyName: "staging-1",
							ClusterURL:   "https://does-not-exist.local",
						},
					},
				},
			},
			DestinationNamespace: "example",
		}

		serviceBytes, err := yaml.Marshal(service)
		if err != nil {
			return fmt.Errorf("failed to marshal service object: %w", err)
		}

		err = os.WriteFile("examples/example-service.yaml", serviceBytes, 0644)
		if err != nil {
			return fmt.Errorf("failed to write service file: %w", err)
		}

		return nil
	}, func(error) {
		cancel()
	})

	g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))

	err = g.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
