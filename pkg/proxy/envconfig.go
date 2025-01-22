package proxy

import (
	"github.com/kelseyhightower/envconfig"
)

// Envconfig is an interface that provides a proxy of the methods of envconfig.
type Envconfig interface {
	Process(prefix string, spec interface{}) error
}

// envconfigProxy is a proxy struct that implements the Envconfig interface.
type envconfigProxy struct{}

// NewEnvconfig returns a new instance of the Envconfig interface.
func NewEnvconfig() Envconfig {
	return &envconfigProxy{}
}

// Process processes the environment variables and stores the result in the spec.
func (envconfigProxy) Process(prefix string, spec interface{}) error {
	return envconfig.Process(prefix, spec)
}
