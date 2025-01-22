package proxy

import (
	"context"
	"time"

	"github.com/eiannone/keyboard"
)

// Keyboard is an interface that provides a proxy of the methods of keyboard.
type Keyboard interface {
	Close()
	GetKey(timeoutSec int) (rune, keyboard.Key, error)
	Open() error
}

// KeyboardProxy is a struct that provides a proxy of the methods of keyboard.
type keyboardProxy struct{}

// NewKeyboard returns a new instance of the keyboard proxy.
func NewKeyboard() Keyboard {
	return &keyboardProxy{}
}

// Open is a proxy method that calls the Open method of the keyboard.
func (k *keyboardProxy) Open() error {
	return keyboard.Open()
}

// Close is a proxy method that calls the Close method of the keyboard.
func (k *keyboardProxy) Close() {
	keyboard.Close()
}

// GetKey is a proxy method that calls the GetKey method of the keyboard.
func (k *keyboardProxy) GetKey(timeoutSec int) (rune, keyboard.Key, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	keyChan := make(chan struct {
		r rune
		k keyboard.Key
		e error
	})

	go func() {
		r, k, e := keyboard.GetKey()
		keyChan <- struct {
			r rune
			k keyboard.Key
			e error
		}{r, k, e}
	}()

	select {
	case result := <-keyChan:
		return result.r, result.k, result.e
	case <-ctx.Done():
		return 0, keyboard.KeyEnter, nil
	}
}
