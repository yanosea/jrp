package userproxy

import (
	"os/user"
)

// UserInstanceInterface is a interface for user.User.
type UserInstanceInterface interface{}

// UserInstance is a struct that implements UserInstanceInterface.
type UserInstance struct {
	FieldUser *user.User
}
