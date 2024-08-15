package usermanager

import (
	"os/user"
)

type UserProvider interface {
	Current() (*user.User, error)
}

type OSUserProvider struct{}

func (OSUserProvider) Current() (*user.User, error) {
	return user.Current()
}
