package goconf

import (
	"fmt"
	"io"
	"os"
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

// tempFilesTxt holds the contents of temporary files which will be created for
// testing purposes. Keys are arbitrary IDs which are used to refer to these
// temporary files in the future.
var tempFilesTxt = map[string][]byte{
	"a": []byte("key1 = \"value1\""),
	"b": []byte("key2 = \"value2\""),
	"c": []byte("key3 = \"value3\""),
}

// configFile is a dummy configuration file struct definition passed to
// Loader.Load in a few tests
type configFile struct {
	Key1 string `mapstructure:"key1" validate:"required"`
	Key2 string `mapstructure:"key2" validate:"required"`
	Key3 string `mapstructure:"key3" validate:"required"`
}

// defaultConfigFile is a dummy configuration file struct compatible with the
// configFile struct, but with a default value specified for Key3
type defaultConfigFile struct {
	Key1 string `mapstructure:"key1" validate:"required"`
	Key2 string `mapstructure:"key2" validate:"required" default:"key2default"`
	Key3 string `mapstructure:"key3" validate:"required" default:"key3default"`
}

// expectedConfigFile holds the values which Loader.Load should place in a
// configFile struct when called
var expectedConfigFile = configFile{
	Key1: "value1",
	Key2: "value2",
	Key3: "value3",
}

// placeTempFiles creates temporary files for testing purposes. The names
// argumen can only include keys which are present in the tempFilesTxt map.
func placeTempFiles(t *testing.T, names []string) map[string]*os.File {
	// Place temp files
	files := map[string]*os.File{}

	for _, name := range names {
		path := fmt.Sprintf("/tmp/goconf-config-%s.toml", name)
		f, err := os.Create(path)

		assert.NilError(t, err, "failed to open config file %s", name)

		files[name] = f
	}

	for name, f := range files {
		_, err := f.Write(tempFilesTxt[name])
		assert.NilError(t, err, "failed to write to config file %s",
			name)
	}

	return files
}

// cleanupTempFiles removes temp files provided by PlaceTempFiles. The files
// argument should be the map returned by placeTempFiles.
func cleanupTempFiles(t *testing.T, files map[string]*os.File) {
	for name, f := range files {
		assert.NilError(t, f.Close(), "failed to close config"+
			"file %s", name)
	}

	for name := range files {
		path := fmt.Sprintf("/tmp/goconf-config-%s.toml", name)
		assert.NilError(t, os.Remove(path), "failed to remove"+
			"config file %s", name)
	}
}

// TestRegisterFormat ensures MapDecoders are properly added to the formats map
func TestRegisterFormat(t *testing.T) {
	// Setup loader
	loader := NewLoader()

	loader.RegisterFormat(".foo", DummyMapDecoder{ID: ".foo"})
	loader.RegisterFormat("", DummyMapDecoder{ID: ""})

	// Assert
	fooDecoder, ok := loader.formats[".foo"]
	assert.Equal(t, ok, true, ".foo dummy map decoder not present")
	assert.Equal(t, fooDecoder.(DummyMapDecoder).ID, ".foo")

	emptyDecoder, ok := loader.formats[""]
	assert.Equal(t, ok, true, "empty string dummy map decoder not present")
	assert.Equal(t, emptyDecoder.(DummyMapDecoder).ID, "")
}

// TestAddConfigPath ensures paths are added to the configPaths array
func TestAddConfigPath(t *testing.T) {
	// Setup loader
	loader := NewLoader()

	loader.AddConfigPath("AAA")
	loader.AddConfigPath("BBB")
	loader.AddConfigPath("CCC")

	// Assert
	assert.DeepEqual(t, loader.configPaths, []string{"AAA", "BBB", "CCC"})
}

// TestLoad ensures Load reads and decodes configuration files
func TestLoad(t *testing.T) {
	// Place temp files
	tmpFiles := placeTempFiles(t, []string{"a", "b", "c"})
	defer cleanupTempFiles(t, tmpFiles)

	// Setup loader
	loader := NewDefaultLoader()

	loader.AddConfigPath("/tmp/goconf-config-*.toml")

	// Load
	actualConfig := configFile{}

	err := loader.Load(&actualConfig)

	// Assert
	assert.NilError(t, err, "failed to load configuration")

	assert.DeepEqual(t, actualConfig, expectedConfigFile)
}

// TestLoadValidate ensures Load validates after reading files into a struct
func TestLoadValidate(t *testing.T) {
	// Place temp files
	tmpFiles := placeTempFiles(t, []string{"a", "b"})
	defer cleanupTempFiles(t, tmpFiles)

	// Setup loader
	loader := NewDefaultLoader()

	loader.AddConfigPath("/tmp/goconf-config-*.toml")

	// Load
	actualConfig := configFile{}

	err := loader.Load(&actualConfig)

	// Assert
	assert.Equal(t, err.Error(), "failed to validate configuration "+
		"struct: Key: 'configFile.Key3' Error:Field validation for "+
		"'Key3' failed on the 'required' tag")
}

// TestLoadDirCheck ensures Load errors if a config path is a directory
func TestLoadDirCheck(t *testing.T) {
	// Setup loader
	loader := NewDefaultLoader()

	loader.AddConfigPath(".")

	// Load
	cfg := configFile{}

	err := loader.Load(&cfg)

	assert.Equal(t, err.Error(), "configuration path \".\" is a directory"+
		", cannot be")
}

// TestLoadDefaults ensures Load sets default values in a struct
func TestLoadDefaults(t *testing.T) {
	// Place temp files
	tmpFiles := placeTempFiles(t, []string{"a", "b"})
	defer cleanupTempFiles(t, tmpFiles)

	// Setup loader
	loader := NewDefaultLoader()

	loader.AddConfigPath("/tmp/goconf-config-*.toml")

	// Load
	cfg := defaultConfigFile{}

	err := loader.Load(&cfg)

	// Assert
	assert.NilError(t, err)

	expectedCfg := defaultConfigFile{
		Key1: expectedConfigFile.Key1,
		Key2: expectedConfigFile.Key2,
		Key3: "key3default",
	}

	assert.DeepEqual(t, cfg, expectedCfg)
}
