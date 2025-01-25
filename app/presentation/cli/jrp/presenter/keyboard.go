package presenter

import (
	"github.com/yanosea/jrp/v2/pkg/proxy"
	"github.com/yanosea/jrp/v2/pkg/utility"
)

var (
	// Ku is a variable that contains the KeyboardUtil struct for injecting dependencies in testing.
	Ku = utility.NewKeyboardUtil(proxy.NewKeyboard())
)

// CloseKeyboard closes the keyboard.
func CloseKeyboard() {
	Ku.CloseKeyboard()
}

// GetKey gets a key from the keyboard.
func GetKey(timeoutSec int) (string, error) {
	return Ku.GetKey(timeoutSec)
}

// OpenKeyboard opens the keyboard.
func OpenKeyboard() error {
	return Ku.OpenKeyboard()
}
