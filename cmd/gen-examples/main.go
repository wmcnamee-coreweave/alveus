package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	"github.com/goccy/go-yaml"
	"github.com/oklog/run"

	"github.com/ghostsquad/alveus/api/v1alpha1"
)

var version string

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var g run.Group
	var err error

	g.Add(func() error {
		return createServiceFile()
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

func createServiceFile() error {
	service := &v1alpha1.Service{
		Name: "example-service",
		Source: v1alpha1.Source{
			Path: "./examples/example-service/manifests",
		},
		IgnoreDifferences:     nil,
		PrePromotionAnalysis:  nil,
		PostPromotionAnalysis: nil,
		DestinationGroups: []v1alpha1.DestinationGroup{
			{
				Name: "staging",
				Destinations: []v1alpha1.Destination{
					{
						ApplicationDestination: &argov1alpha1.ApplicationDestination{
							Server: "http://kube.local",
						},
						ArgoCDLogin: v1alpha1.ArgoCDLogin{},
					},
				},
			},
		},
		DestinationNamespace: "my-namespace",
		SyncPolicy:           nil,
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
}
