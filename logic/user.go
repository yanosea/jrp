package logic

import (
	"os/user"
)

type User interface {
	Current() (*user.User, error)
}

type OsUser struct{}

func (u OsUser) Current() (*user.User, error) {
	return user.Current()
}
