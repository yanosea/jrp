package sortproxy

import (
	"sort"
)

// Sort is an interface for sort.
type Sort interface {
	Slice(x any, less func(i, j int) bool)
}

// SortProxy is a struct that implements Sort.
type SortProxy struct{}

// New is a constructor of SortProxy.
func New() Sort {
	return &SortProxy{}
}

// Slice is a proxy for sort.Slice.
func (*SortProxy) Slice(x any, less func(i, j int) bool) {
	sort.Slice(x, less)
}
