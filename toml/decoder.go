package toml

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

// TomlMapDecoder implements MapDecoder for Toml files
type TomlMapDecoder struct{}

// Decode Toml file into map
func (d TomlMapDecoder) Decode(r io.Reader, m *map[string]interface{}) error {
	if _, err := toml.DecodeReader(r, m); err != nil {
		return fmt.Errorf("error decoding toml: %s", err.Error())
	}

	return nil
}
