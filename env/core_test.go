package env

import (
	"testing"

	"github.com/Bios-Marcel/yagcl"
	"github.com/stretchr/testify/assert"
)

func Test_EventSource_InterfaceCompliance(t *testing.T) {
	var _ yagcl.Source = Source()
}

func Test_Parse_KeyTags(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
		FieldB string `env:"FIELD_B"`
	}

	defer setEnvTemporarily("FIELD_A", "content a")()
	defer setEnvTemporarily("FIELD_B", "content b")()
	var c configuration
	err := yagcl.New[configuration]().
		Add(Source()).
		Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "content a", c.FieldA)
		assert.Equal(t, "content b", c.FieldB)
	}
}

func Test_Parse_Prefix(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
		FieldB string `env:"FIELD_B"`
	}

	defer setEnvTemporarily("TEST_FIELD_A", "content a")()
	defer setEnvTemporarily("TEST_FIELD_B", "content b")()
	var c configuration
	err := yagcl.
		New[configuration]().
		Add(Source().Prefix("TEST")).
		Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "content a", c.FieldA)
		assert.Equal(t, "content b", c.FieldB)
	}
}
func Test_Parse_KeyValueConverter(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
		FieldB string `env:"FIELD_B"`
	}

	defer setEnvTemporarily("TEST_field_a", "content a")()
	defer setEnvTemporarily("TEST_FIELD_B", "content b")()
	var c configuration
	err := yagcl.
		New[configuration]().
		Add(
			Source().
				Prefix("TEST_").
				KeyValueConverter(func(s string) string {
					return s
				}),
		).
		Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "content a", c.FieldA)
		assert.Equal(t, "content b", c.FieldB)
	}
}

func Test_Parse_MissingFieldKey(t *testing.T) {
	type configuration struct {
		FieldA string
	}

	var c configuration
	err := yagcl.
		New[configuration]().
		Add(Source()).
		Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrExportedFieldMissingKey)
}

func Test_Parse_IgnoreField(t *testing.T) {
	type configuration struct {
		FieldA string `ignore:"true"`
	}

	var c configuration
	err := yagcl.
		New[configuration]().
		Add(Source()).
		Parse(&c)
	if assert.NoError(t, err) {
		assert.Empty(t, c.FieldA)
	}
}
