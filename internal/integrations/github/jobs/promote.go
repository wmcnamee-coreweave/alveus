package jobs

import (
	gocto "github.com/cakehappens/gocto"
)

func NewPromotion(workflowName, workflowFilename string) gocto.Job {
	job := gocto.Job{
		Name: "promote-to-" + workflowName,
		Uses: ".github/workflows/" + workflowFilename,
		With: nil,
	}

	return job
}
