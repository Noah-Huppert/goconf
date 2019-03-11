package main

// Loader loads configuration
type Loader struct {
	// formats holds the MapDecoders for file extensions
	formats map[string]MapDecoder

	// configPaths are paths to files to load
	configPaths []string
}
