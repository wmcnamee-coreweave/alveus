package workflows

import (
	gocto "github.com/cakehappens/gocto"
)

func NewDeployment(name string, stages ...gocto.Job) gocto.Workflow {
	wf := gocto.Workflow{
		Name: name,
		On:   gocto.WorkflowOn{},
		Jobs: make(map[string]gocto.Job),
	}

	for _, stage := range stages {
		wf.Jobs[stage.Name] = stage
	}

	return wf
}
