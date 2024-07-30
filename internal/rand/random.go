package rand

import (
	"math/rand"
	"time"
)

type RandomGenerator interface {
	Intn(n int) int
}

type DefaultRandomGenerator struct {
	rng *rand.Rand
}

func NewDefaultRandomGenerator() *DefaultRandomGenerator {
	return &DefaultRandomGenerator{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (d *DefaultRandomGenerator) Intn(n int) int {
	return d.rng.Intn(n)
}
