package models

import "github.com/ghostsquad/alveus/api/v1alpha1"

type ConfigFile struct {
	Path   string
	Config v1alpha1.Config
}
