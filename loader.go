package goconf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Noah-Huppert/goconf/toml"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mcuadros/go-defaults.v1"
)

// Loader loads configuration
type Loader struct {
	// validate is the validator instance to use for validation
	validate *validator.Validate

	// formats holds the MapDecoders for file extensions
	formats map[string]MapDecoder

	// configPaths are paths to files to load
	configPaths []string
}

// NewLoader creates a Loader
func NewLoader() *Loader {
	return &Loader{
		validate:    validator.New(),
		formats:     map[string]MapDecoder{},
		configPaths: []string{},
	}
}

// NewDefaultLoader creates a loader with all MapDecoders implemented by
// goconf registered with the appropriate file extensions.
func NewDefaultLoader() *Loader {
	l := NewLoader()

	l.RegisterFormat(".toml", toml.TomlMapDecoder{})

	return l
}

// RegisterFormat registers a MapDecoder to be used for a file extension.
// The ext argument should include the final dot and then the extension name.
// An empty string can be passed to target files without an extension.
func (l *Loader) RegisterFormat(ext string, decoder MapDecoder) {
	l.formats[ext] = decoder
}

// GetValidate returns the validator instance used to validate configuration
func (l Loader) GetValidate() *validator.Validate {
	return l.validate
}

// AddConfigPath adds a potential path from which configuration files will be
// loaded. Must point to file(s) not of directories. The p argument can contain
// shell globs.
func (l *Loader) AddConfigPath(p string) {
	// {{{1 Check if already in configPaths
	for _, existingPath := range l.configPaths {
		if existingPath == p {
			return
		}
	}

	// {{{1 Add
	l.configPaths = append(l.configPaths, p)
}

// Load configuration files into a struct
func (l Loader) Load(c interface{}) error {
	// {{{1 Set default values
	defaults.SetDefaults(c)

	// {{{1 Expand config paths
	loadPaths := []string{}

	for _, configPath := range l.configPaths {
		// {{{2 Interpret shell globs
		expandedPaths, err := filepath.Glob(configPath)
		if err != nil {
			return fmt.Errorf("failed to expand configuration "+
				"path \"%s\" glob: %s", configPath,
				err.Error())
		}

		for _, expandedPath := range expandedPaths {
			// {{{2 Check not directory
			fi, err := os.Stat(expandedPath)
			if err != nil {
				return fmt.Errorf("failed to stat "+
					"configuration path \"%s\": %s",
					expandedPath, err.Error())
			}

			if fi.IsDir() {
				return fmt.Errorf("configuration path "+
					"\"%s\" is a directory, cannot be",
					expandedPath)
			}

			// {{{2 Not directory, add
			loadPaths = append(loadPaths, expandedPath)
		}
	}

	// {{{1 Try to load all files in loadPaths
	for _, loadPath := range loadPaths {
		// {{{2 Check if MapDecoder exists for file extension
		decoder, ok := l.formats[filepath.Ext(loadPath)]

		if !ok {
			continue
		}

		// {{{2 Use MapDecoder if exists
		// {{{3 Open file
		loadFile, err := os.Open(loadPath)
		if err != nil {
			return fmt.Errorf("error opening configuration "+
				"file \"%s\": %s", loadPath, err.Error())
		}

		// {{{3 Call MapDecoder
		loadMap := map[string]interface{}{}

		if err = decoder.Decode(loadFile, &loadMap); err != nil {
			return fmt.Errorf("error decoding \"%s\": %s",
				loadPath, err.Error())
		}

		// {{{3 Put map into struct
		if err = mapstructure.Decode(loadMap, c); err != nil {
			return fmt.Errorf("error putting decoded map "+
				"for \"%s\" into configuration struct: %s",
				loadPath, err.Error())
		}
	}

	// {{{1 Validate configuration struct
	if err := l.validate.Struct(c); err != nil {
		return fmt.Errorf("failed to validate configuration "+
			"struct: %s", err.Error())
	}

	return nil
}
