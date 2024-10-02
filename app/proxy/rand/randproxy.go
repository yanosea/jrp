package randproxy

import (
	"math/rand"
)

// Rand is an interface for rand.
type Rand interface {
	Intn(n int) int
}

// RandProxy is a struct that implements Rand.
type RandProxy struct{}

// New is a constructor for RandProxy.
func New() Rand {
	return &RandProxy{}
}

// Intn is a proxy for rand.Intn.
func (*RandProxy) Intn(n int) int {
	return rand.Intn(n)
}
