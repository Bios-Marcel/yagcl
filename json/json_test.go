package json

import (
	"testing"

	"github.com/Bios-Marcel/yagcl"
	"github.com/stretchr/testify/assert"
)

func Test_JSONSource_InterfaceCompliance(t *testing.T) {
	var _ yagcl.Source = Source("")
}

func Test_Parse_JSON_Simple(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
		FieldB string `json:"field_b"`
	}
	var c configuration
	err := yagcl.New[configuration]().
		AddSource(Source("./test.json").Must()).
		Parse(&c)
	assert.NoError(t, err)
	assert.Equal(t, "content a", c.FieldA)
	assert.Equal(t, "content b", c.FieldB)
}
