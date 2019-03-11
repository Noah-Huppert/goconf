package main

import (
	"io"
)

// MapDecoder decodes a Reader into a map
type MapDecoder interface {
	// Decode a reader into a map
	Decode(r io.Reader, m *map[string]interface{}) error
}
