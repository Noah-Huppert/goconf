package goconf_test

import (
	"github.com/Noah-Huppert/goconf"
	"github.com/Noah-Huppert/goconf/toml"
)

// Simple configuration setup using Toml files
func Example_toml() {
	// import "github.com/Noah-Huppert/goconf"
	// import "github.com/Noah-Huppert/goconf/toml"

	// Create goconf instance
	loader := goconf.NewLoader()

	// Register file formats
	loader.RegisterFormat(".toml", toml.TomlMapDecoder{})

	// Define locations to search for configuration files
	// Can use shell globs
	loader.AddConfigPath("/etc/foo/foo.*")
	loader.AddConfigPath("/etc/foo.d/*")

	// Load values
	type YourConfigStruct struct {
		Foo string `mapstructure:"foo"`
		Bar string `mapstructure:"bar"`
	}

	config := YourConfigStruct{}
	err := loader.Load(&config)
	panic(err)
}
