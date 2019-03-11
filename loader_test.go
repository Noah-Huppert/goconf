package goconf

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/Noah-Huppert/goconf/toml"

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

var configFilesTxt map[string][]byte = map[string][]byte{
	"a": []byte("key1 = \"value1\""),
	"b": []byte("key2 = \"value2\""),
	"c": []byte("key3 = \"value3\""),
}

type configFile struct {
	Key1 string `mapstructure:"key1"`
	Key2 string `mapstructure:"key2"`
	Key3 string `mapstructure:"key3"`
}

var expectedConfigFile configFile = configFile{
	Key1: "value1",
	Key2: "value2",
	Key3: "value3",
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

// TestAddConfigPath ensures paths are added to the configPaths array
func TestAddConfigPath(t *testing.T) {
	loader := NewLoader()

	loader.AddConfigPath("AAA")
	loader.AddConfigPath("BBB")
	loader.AddConfigPath("CCC")

	assert.DeepEqual(t, loader.configPaths, []string{"AAA", "BBB", "CCC"})
}

// TestLoad ensures Load reads and decodes configuration files
func TestLoad(t *testing.T) {
	// Place temp files
	files := map[string]*os.File{}

	for _, name := range []string{"a", "b", "c"} {
		path := fmt.Sprintf("/tmp/goconf-config-%s.toml", name)
		f, err := os.Create(path)

		assert.NilError(t, err, "failed to open config file %s", name)

		files[name] = f
	}

	for name, f := range files {
		_, err := f.Write(configFilesTxt[name])
		assert.NilError(t, err, "failed to write to config file %s",
			name)
	}

	defer func() {
		for name, f := range files {
			assert.NilError(t, f.Close(), "failed to close config"+
				"file %s", name)
		}
	}()

	defer func() {
		for name, _ := range files {
			path := fmt.Sprintf("/tmp/goconf-config-%s.toml", name)
			assert.NilError(t, os.Remove(path), "failed to remove"+
				"config file %s", name)
		}
	}()

	// Setup loader
	loader := NewLoader()

	loader.AddConfigPath("/tmp/goconf-config-*.toml")
	loader.RegisterFormat(".toml", toml.TomlMapDecoder{})

	// Load
	actualConfig := configFile{}

	err := loader.Load(&actualConfig)

	assert.NilError(t, err, "failed to load configuration")

	assert.DeepEqual(t, actualConfig, expectedConfigFile)
}
