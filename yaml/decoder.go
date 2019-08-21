package yaml

import (
	"io"

	yamlLib "gopkg.in/yaml.v2"
)

// YamlMapDecoder decodes YAML files into maps
type YamlMapDecoder struct{}

// Decode Yaml file into map
func (d YamlMapDecoder) Decode(r io.Reader, m *map[string]interface{}) error {
	decoder := yamlLib.NewDecoder(r)
	return decoder.Decode(m)
}
