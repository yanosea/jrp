package keyboardproxy

import (
	"context"
	"time"

	"github.com/eiannone/keyboard"
)

// Keyboard is an interface for keyboard.
type Keyboard interface {
	Open() error
	Close()
	GetKey(timeoutSec int) (rune, keyboard.Key, error)
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
func (*KeyboardProxy) GetKey(timeoutSec int) (rune, keyboard.Key, error) {
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
