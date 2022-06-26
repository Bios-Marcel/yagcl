package env

import (
	"fmt"
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

func Test_Parse_KeyTags(t *testing.T) {
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
		AddSource(Source().Prefix("TEST")).
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
		AddSource(
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

func Test_Parse_String_Valid(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "content a")()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "content a", c.FieldA)
	}
}
func Test_Parse_Struct_Valid(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
		FieldB struct {
			FieldC string `key:"field_c"`
		} `key:"field_b"`
	}

	defer setEnvTemporarily("FIELD_A", "content a")()
	defer setEnvTemporarily("FIELD_B_FIELD_C", "content c")()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "content a", c.FieldA)
		assert.Equal(t, "content c", c.FieldB.FieldC)
	}
}

func Test_Parse_String_Whitespace(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "   ")()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "   ", c.FieldA)
	}
}

const maxUint = ^uint(0)
const minUint = uint(0)
const maxInt = int(maxUint >> 1)
const minInt = -maxInt - 1

const maxUint8 = ^uint8(0)
const minUint8 = uint8(0)
const maxInt8 = int8(maxUint8 >> 1)
const minInt8 = -maxInt8 - 1

const maxUint16 = ^uint16(0)
const minUint16 = uint16(0)
const maxInt16 = int16(maxUint16 >> 1)
const minInt16 = -maxInt16 - 1

const maxUint32 = ^uint32(0)
const minUint32 = uint32(0)
const maxInt32 = int32(maxUint32 >> 1)
const minInt32 = -maxInt32 - 1

const maxUint64 = ^uint64(0)
const minUint64 = uint64(0)
const maxInt64 = int64(maxUint64 >> 1)
const minInt64 = -maxInt64 - 1

func Test_Parse_Int_Valid(t *testing.T) {
	type configuration struct {
		Min int `key:"min"`
		Max int `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", minInt))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", maxInt))()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, minInt, c.Min)
		assert.Equal(t, maxInt, c.Max)
	}
}

func Test_Parse_Int8_Valid(t *testing.T) {
	type configuration struct {
		Min int8 `key:"min"`
		Max int8 `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", minInt8))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", maxInt8))()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, minInt8, c.Min)
		assert.Equal(t, maxInt8, c.Max)
	}
}

func Test_Parse_Int16_Valid(t *testing.T) {
	type configuration struct {
		Min int16 `key:"min"`
		Max int16 `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", minInt16))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", maxInt16))()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, minInt16, c.Min)
		assert.Equal(t, maxInt16, c.Max)
	}
}

func Test_Parse_Int32_Valid(t *testing.T) {
	type configuration struct {
		Min int32 `key:"min"`
		Max int32 `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", minInt32))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", maxInt32))()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, minInt32, c.Min)
		assert.Equal(t, maxInt32, c.Max)
	}
}

func Test_Parse_Int64_Valid(t *testing.T) {
	type configuration struct {
		Min int64 `key:"min"`
		Max int64 `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", minInt64))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", maxInt64))()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, minInt64, c.Min)
		assert.Equal(t, maxInt64, c.Max)
	}
}

func Test_Parse_Uint_Valid(t *testing.T) {
	type configuration struct {
		Min uint `key:"min"`
		Max uint `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", minUint))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", maxUint))()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, minUint, c.Min)
		assert.Equal(t, maxUint, c.Max)
	}
}

func Test_Parse_Uint8_Valid(t *testing.T) {
	type configuration struct {
		Min uint8 `key:"min"`
		Max uint8 `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", minUint8))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", maxUint8))()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, minUint8, c.Min)
		assert.Equal(t, maxUint8, c.Max)
	}
}

func Test_Parse_Uint16_Valid(t *testing.T) {
	type configuration struct {
		Min uint16 `key:"min"`
		Max uint16 `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", minUint16))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", maxUint16))()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, minUint16, c.Min)
		assert.Equal(t, maxUint16, c.Max)
	}
}

func Test_Parse_Uint32_Valid(t *testing.T) {
	type configuration struct {
		Min uint32 `key:"min"`
		Max uint32 `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", minUint32))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", maxUint32))()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, minUint32, c.Min)
		assert.Equal(t, maxUint32, c.Max)
	}
}

func Test_Parse_Uint64_Valid(t *testing.T) {
	type configuration struct {
		Min uint64 `key:"min"`
		Max uint64 `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", minUint64))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", maxUint64))()
	var c configuration
	err := yagcl.New[configuration]().AddSource(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, minUint64, c.Min)
		assert.Equal(t, maxUint64, c.Max)
	}
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

func Test_Parse_Uint_Invalid(t *testing.T) {
	type configuration struct {
		FieldA uint `key:"field_a"`
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
	if assert.NoError(t, err) {
		assert.Equal(t, "i am the default", c.FieldA)
	}
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
	if assert.NoError(t, err) {
		assert.Empty(t, c.FieldA)
	}
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
