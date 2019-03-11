package toml

import (
	"io"
	"strings"
	"testing"

	"gotest.tools/assert"
)

var testTomlReader io.Reader = strings.NewReader(`
key1 = "value1"
key2 = "value2"

[table1]
key3 = "value3"
`)

var testTomlValue map[string]interface{} = map[string]interface{}{
	"key1": "value1",
	"key2": "value2",
	"table1": map[string]interface{}{
		"key3": "value3",
	},
}

func TestDecode(t *testing.T) {
	decoder := TomlMapDecoder{}
	m := map[string]interface{}{}

	err := decoder.Decode(testTomlReader, &m)
	assert.NilError(t, err, "failed to decode toml")

	assert.DeepEqual(t, testTomlValue, m)
}
