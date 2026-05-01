package platform

import "github.com/atotto/clipboard"

// Clipboard defines the interface for clipboard operations (Strategy Pattern).
type Clipboard interface {
	Write(text string) error
}

// SystemClipboard implements Clipboard using the system's native clipboard.
type SystemClipboard struct{}

func (c *SystemClipboard) Write(text string) error {
	return clipboard.WriteAll(text)
}
