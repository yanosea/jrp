package proxy

import (
	"math/rand"
)

// Rand is an interface that provides a proxy of the methods of math/rand.
type Rand interface {
	Intn(n int) int
}

// randProxy is a proxy struct that implements the Rand interface.
type randProxy struct{}

// NewRand returns a new instance of the Rand interface.
func NewRand() Rand {
	return &randProxy{}
}

// Intn returns, as an int, a non-negative pseudo-random number in [0,n).
func (r *randProxy) Intn(n int) int {
	return rand.Intn(n)
}
