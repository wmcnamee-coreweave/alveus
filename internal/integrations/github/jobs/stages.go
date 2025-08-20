package jobs

import (
	gocto "github.com/cakehappens/gocto"
	"github.com/lithammer/dedent"
)

type Options struct {
	// Environment defaults to the name of the jobs
	Environment gocto.Environment
}

type Option func(opts *Options)

func WithEnvironment(env gocto.Environment) Option {
	return func(opts *Options) {
		opts.Environment = env
	}
}

func WithEnvironmentURL(url string) Option {
	return func(opts *Options) {
		opts.Environment.URL = url
	}
}

func WithEnvironmentName(name string) Option {
	return func(opts *Options) {
		opts.Environment.Name = name
	}
}

func NewStage(name string, options ...Option) gocto.Job {
	opts := &Options{
		Environment: gocto.Environment{
			Name: name,
		},
	}
	for _, option := range options {
		option(opts)
	}

	job := gocto.Job{

		Name:        "jobs-" + name,
		Permissions: gocto.Permissions{},
		Environment: opts.Environment,
		Concurrency: gocto.Concurrency{
			Group:            "jobs-" + name,
			CancelInProgress: false,
		},
		Defaults: gocto.Defaults{
			Run: gocto.DefaultsRun{
				Shell: gocto.ShellBash,
			},
		},
		Steps: []gocto.Step{
			{
				Uses: "actions/checkout@v4",
			},
			{
				Run: dedent.Dedent(`
							argocd login 
						`),
			},
			{
				Run: dedent.Dedent(`
							argocd app create --upsert
						`),
			},
		},
	}

	return job
}
