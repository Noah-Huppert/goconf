package main

import (
	"io"
	"testing"

	"gotest.tools/assert"
)

// DummyMapDecoder is used to test the RegisterFormat method
type DummyMapDecoder struct {
	// ID identifies the dummy decoder for test purposes
	ID string
}

// Decode implements MapDecoder for DummyMapDecoder
func (d DummyMapDecoder) Decode(r io.Reader, m *map[string]interface{}) error {
	return nil
}

// TestRegisterFormat ensures MapDecoders are properly added to the formats map
func TestRegisterFormat(t *testing.T) {
	loader := NewLoader()

	loader.RegisterFormat(".foo", DummyMapDecoder{ID: ".foo"})
	loader.RegisterFormat("", DummyMapDecoder{ID: ""})

	fooDecoder, ok := loader.formats[".foo"]
	assert.Equal(t, ok, true, ".foo dummy map decoder not present")
	assert.Equal(t, fooDecoder.(DummyMapDecoder).ID, ".foo")

	emptyDecoder, ok := loader.formats[""]
	assert.Equal(t, ok, true, "empty string dummy map decoder not present")
	assert.Equal(t, emptyDecoder.(DummyMapDecoder).ID, "")
}
