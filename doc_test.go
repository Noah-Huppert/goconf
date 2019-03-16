package goconf_test

import (
	"github.com/Noah-Huppert/goconf"
)

// Simple configuration setup using Toml files
func Example_toml() {
	// import "github.com/Noah-Huppert/goconf"

	// Create goconf instance
	loader := goconf.NewDefaultLoader()

	// Define locations to search for configuration files
	// Can use shell globs
	loader.AddConfigPath("/etc/foo/foo.*")
	loader.AddConfigPath("/etc/foo.d/*")

	// Load values
	type YourConfigStruct struct {
		Foo string `mapstructure:"foo" validate:"required"`
		Bar string `mapstructure:"bar"`
	}

	config := YourConfigStruct{}
	err := loader.Load(&config)
	panic(err)
}
