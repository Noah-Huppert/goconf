package main

// Loader loads configuration
type Loader struct {
	// formats holds the MapDecoders for file extensions
	formats map[string]MapDecoder

	// configPaths are paths to files to load
	configPaths []string
}

// NewLoader creates a Loader
func NewLoader() {
	return Loader{
		formats:     map[string]MapDecoder{},
		configPaths: []string{},
	}
}

// RegisterFormat adds a MapDecoder to the formats field
func (l *Loader) RegisterFormat(ext string, decoder MapDecoder) {
	l.formats[ext] = decoder
}

// AddConfigPath adds a path to the configPaths field
func (l *Loader) AddConfigPath(p string) {
	// Check if already in configPaths
	for existingPath := range l.configPaths {
		if existingPath == p {
			return
		}
	}

	// Add
	l.configPaths = append(l.configPaths, p)
}
