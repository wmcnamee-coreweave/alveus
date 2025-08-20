package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/oklog/run"

	"github.com/ghostsquad/alveus/internal/cmd"
)

var version string

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var g run.Group
	var err error

	rootCmd := cmd.NewRootCommand()
	rootCmd.SetContext(ctx)

	g.Add(func() error {
		return rootCmd.Execute()
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
