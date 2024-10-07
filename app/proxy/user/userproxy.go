package userproxy

import (
	"os/user"
)

// User is an interface for user.
type User interface {
	Current() (*UserInstance, error)
}

// UserProxy is a struct that implements User.
type UserProxy struct{}

// New is a constructor for UserProxy.
func New() User {
	return &UserProxy{}
}

// Current is a proxy for user.Current.
func (*UserProxy) Current() (*UserInstance, error) {
	currentUser, err := user.Current()
	return &UserInstance{FieldUser: currentUser}, err
}
