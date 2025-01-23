package utility

import (
	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// RandUtil is an interface that contains the utility functions for generating random numbers.
type RandUtil interface {
	GenerateRandomNumber(max int) int
}

// randUtil is a struct that contains the utility functions for generating random numbers.
type randUtil struct {
	rand proxy.Rand
}

// NewRandUtil returns a new instance of the RandomUtil struct.
func NewRandUtil(rand proxy.Rand) RandUtil {
	return &randUtil{
		rand: rand,
	}
}

// GenerateRandomNumber generates a random number between min and max.
func (ru *randUtil) GenerateRandomNumber(max int) int {
	if max <= 0 {
		return 0
	}
	return ru.rand.Intn(max)
}
