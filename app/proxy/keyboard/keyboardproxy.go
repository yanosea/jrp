package keyboardproxy

import (
	"github.com/eiannone/keyboard"
)

// Keyboard is an interface for keyboard.
type Keyboard interface {
	Open() error
	Close()
	GetKey() (rune, keyboard.Key, error)
}

// KeyboardProxy is a struct that implements Keyboard.
type KeyboardProxy struct{}

// New is a constructor for KeyboardProxy.
func New() Keyboard {
	return &KeyboardProxy{}
}

// Open is a proxy for keyboard.Open.
func (*KeyboardProxy) Open() error {
	return keyboard.Open()
}

// Close is a proxy for keyboard.Close.
func (*KeyboardProxy) Close() {
	keyboard.Close()
}

// GetKey is a proxy for keyboard.GetKey.
func (*KeyboardProxy) GetKey() (rune, keyboard.Key, error) {
	return keyboard.GetKey()
}
