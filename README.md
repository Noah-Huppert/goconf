# Goconf
Simple go configuration library.

# Table Of Contents
- [Overview](#overview)
- [Usage](#usage)

# Overview
Goconf is a simple, straightforward, go configuration library.  

Configuration is defined via structs and tags. Values are loaded from files.  

Any file format can be used. 

# Usage
## Define Configuration
Define configuration parameters in a struct.  

Use [`mapstructure` tags](https://godoc.org/github.com/mitchellh/mapstructure#example-Decode--Tags)
to specify the names of fields when being unmarshalled.  

Use [`validate` tags](https://godoc.org/gopkg.in/go-playground/validator.v9) to
specify value requirements for fields.

## Load values
```go
// Create goconf instance
loader := goconf.NewLoader()

// Register file formats
loader.RegisterFormat(".json", 

// Define locations to search for configuration files
// Can use shell globs
loader.AddConfigPath("/etc/foo/foo.*")
loader.AddConfigPath("/etc/foo.d/*")

// Load values
config := YourConfigStruct{}
err := loader.Load(&config)
```

## Custom File Formats
The `MapUnmarshaller` interface allows Goconf to use any file format.  

Goconf provides an implementation for TOML files in the 
`github.com/Noah-Huppert/goconf-toml` package.

To use any other file format simply implement a `MapUnmarshaller` and register
it with Goconf via the
`Loader.RegisterFormat(string fileExt, MapUnmarshaller unmarshaller)` method.

