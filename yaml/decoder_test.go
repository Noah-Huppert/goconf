package yaml

import (
	"io"
	"strings"
	"testing"

	"gotest.tools/assert"
)

var testYamlReader io.Reader = strings.NewReader(`
key1: value1
key2: value2
table1:
  key3: value3
`)

var expectedYaml map[string]interface{} = map[string]interface{}{
	"key1": "value1",
	"key2": "value2",
	"table1": map[string]interface{}{
		"key3": "value3",
	},
}

// TestDecode ensures YamlMapDecoder properly decodes Yaml
func TestDecode(t *testing.T) {
	decoder := YamlMapDecoder{}
	actualYaml := map[string]interface{}{}

	err := decoder.Decode(testYamlReader, &actualYaml)
	assert.NilError(t, err, "failed to decode yaml")

	assert.DeepEqual(t, expectedYaml, actualYaml)
}
