package utility

import (
	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// KeyboardUtil provides the utility for the keyboard.
type KeyboardUtil interface {
	CloseKeyboard() error
	GetKey(timeoutSec int) (string, error)
	OpenKeyboard() error
}

// keyboardUtil is a struct that implements the KeyboardUtil interface.
type keyboardUtil struct {
	keyboard proxy.Keyboard
}

// NewKeyboardUtil returns a new instance of the KeyboardUtil.
func NewKeyboardUtil(
	keyboard proxy.Keyboard,
) KeyboardUtil {
	return &keyboardUtil{
		keyboard: keyboard,
	}
}

// CloseKeyboard closes the keyboard.
func (k *keyboardUtil) CloseKeyboard() error {
	return k.keyboard.Close()
}

// GetKey returns a key from the keyboard.
func (k *keyboardUtil) GetKey(timeoutSec int) (string, error) {
	rune, _, err := k.keyboard.GetKey(timeoutSec)
	if err != nil {
		return "", err
	}

	return string(rune), nil
}

// OpenKeyboard opens the keyboard.
func (k *keyboardUtil) OpenKeyboard() error {
	return k.keyboard.Open()
}
