package env

import (
	"os"
	"testing"

	"github.com/Bios-Marcel/yagcl"
	"github.com/stretchr/testify/assert"
)

// setEnvTemporarily sets an environment variable and returns a function that
// can reset it to its previous state, whether that state was "set but empty"
// or "unset". This is necessary in order make sure that follow-up tests aren't
// affected by side effects.
func setEnvTemporarily(key, value string) func() {
	oldValue, set := os.LookupEnv(key)
	os.Setenv(key, value)
	if set {
		return func() {
			os.Setenv(key, oldValue)
		}
	}
	return func() {
		os.Unsetenv(key)
	}
}

func Test_EventSource_InterfaceCompliance(t *testing.T) {
	var _ yagcl.Source = Source()
}

func Test_Parse_TestKeyTags(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
		FieldB string `env:"FIELD_B"`
	}

	defer setEnvTemporarily("FIELD_A", "content a")()
	defer setEnvTemporarily("FIELD_B", "content b")()
	var c configuration
	err := yagcl.New[configuration]().
		AddSource(Source()).
		Parse(&c)
	assert.NoError(t, err)
	assert.Equal(t, "content a", c.FieldA)
	assert.Equal(t, "content b", c.FieldB)
}

func Test_Parse_TestPrefix(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
		FieldB string `env:"FIELD_B"`
	}

	defer setEnvTemporarily("TEST_FIELD_A", "content a")()
	defer setEnvTemporarily("TEST_FIELD_B", "content b")()
	var c configuration
	err := yagcl.
		New[configuration]().
		AddSource(Source().Prefix("TEST_")).
		Parse(&c)
	assert.NoError(t, err)
	assert.Equal(t, "content a", c.FieldA)
	assert.Equal(t, "content b", c.FieldB)
}

func Test_Parse_String_Valid(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "content a")()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	assert.NoError(t, err)
	assert.Equal(t, "content a", c.FieldA)
}

func Test_Parse_String_Whitespace(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "   ")()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	assert.NoError(t, err)
	assert.Equal(t, "   ", c.FieldA)
}

func Test_Parse_Int_Valid(t *testing.T) {
	type configuration struct {
		FieldA int `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "1")()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	assert.NoError(t, err)
	assert.Equal(t, 1, c.FieldA)
}

func Test_Parse_Int_Invalid(t *testing.T) {
	type configuration struct {
		FieldA int `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "10no int here")()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrParseValue)
}

func Test_Parse_DefaultValue_String(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a" default:"i am the default"`
	}

	var c configuration
	err := yagcl.
		New[configuration]().
		AddSource(Source()).
		Parse(&c)
	assert.NoError(t, err)
	assert.Equal(t, "i am the default", c.FieldA)
}

func Test_Parse_MissingFieldKey(t *testing.T) {
	type configuration struct {
		FieldA string
	}

	var c configuration
	err := yagcl.
		New[configuration]().
		AddSource(Source()).
		Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrExportedFieldMissingKey)
}

func Test_Parse_IgnoreField1(t *testing.T) {
	type configuration struct {
		FieldA string `ignore:"true"`
	}

	var c configuration
	err := yagcl.
		New[configuration]().
		AddSource(Source()).
		Parse(&c)
	assert.NoError(t, err)
	assert.Empty(t, c.FieldA)
}

func Test_Parse_RequiredValue_Missing_String(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a" required:"true"`
	}

	var c configuration
	err := yagcl.
		New[configuration]().
		AddSource(Source()).
		Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrValueNotSet)
}

func Test_Parse_RequiredValue_EmptyDefault_String(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a" required:"true" default:""`
	}

	var c configuration
	err := yagcl.
		New[configuration]().
		AddSource(Source()).
		Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrValueNotSet)
}
