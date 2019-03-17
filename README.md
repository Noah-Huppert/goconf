# Goconf [![Go Doc](https://godoc.org/github.com/Noah-Huppert/goconf?status.svg)](https://godoc.org/github.com/Noah-Huppert/goconf)
Simple go configuration library.

# Table Of Contents
- [Overview](#overview)
- [Usage](#usage)
- [Tests](#tests)

# Overview
Goconf is a simple, straightforward, go configuration library.  

Configuration is defined via structs and tags. Values are loaded from files.  

Any file format can be used. 

# Usage
## Example
The following example loads TOML configuration files from `/etc/foo/` 
and `/etc/foo.d/`. The configuration files are required to have a `foo` key but 
the `bar` key is optional`.

```go
import "github.com/Noah-Huppert/goconf"

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
```

## Define Configuration
Define configuration parameters in a struct.  

Use [`mapstructure` tags](https://godoc.org/github.com/mitchellh/mapstructure#example-Decode--Tags)
to specify the names of fields when being decoded.  

Use [`validate` tags](https://godoc.org/gopkg.in/go-playground/validator.v9) to
specify value requirements for fields.

## Custom File Formats
The [`MapDecoder`](https://godoc.org/github.com/Noah-Huppert/goconf#MapDecoder)
interface allows Goconf to use any file format.  

Goconf provides an implementation for TOML files in the 
[`github.com/Noah-Huppert/goconf/toml`](https://godoc.org/github.com/Noah-Huppert/goconf/toml) package.

To use any other file format simply implement a [`MapDecoder`](https://godoc.org/github.com/Noah-Huppert/goconf#MapDecoder) 
and register it with Goconf via the
[`Loader.RegisterFormat()`](https://godoc.org/github.com/Noah-Huppert/goconf#Loader.RegisterFormat) method.

# Tests
Run tests:

```
make test
# Or
make
```
