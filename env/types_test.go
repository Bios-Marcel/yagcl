package env

import (
	"encoding/json"
	"fmt"
	"math"
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

func Test_Parse_String_Valid(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "content a")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
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
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "content a", c.FieldA)
		assert.Equal(t, "content c", c.FieldB.FieldC)
	}
}

func Test_Parse_SimplePointer(t *testing.T) {
	type configuration struct {
		FieldA *uint `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "10")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, uint(10), *c.FieldA)
	}
}

func Test_Parse_DoublePointer(t *testing.T) {
	type configuration struct {
		FieldA **uint `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "10")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, uint(10), **c.FieldA)
	}
}

func Test_Parse_PointerOfDoom(t *testing.T) {
	type configuration struct {
		FieldA ***************************************uint `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "10")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, uint(10), ***************************************c.FieldA)
	}
}

func Test_Parse_SinglePointerToStruct(t *testing.T) {
	type substruct struct {
		FieldC string `key:"field_c"`
	}
	type configuration struct {
		FieldA string     `key:"field_a"`
		FieldB *substruct `key:"field_b"`
	}

	defer setEnvTemporarily("FIELD_A", "content a")()
	defer setEnvTemporarily("FIELD_B_FIELD_C", "content c")()
	var c configuration
	c.FieldB = &substruct{}
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "content a", c.FieldA)
		assert.Equal(t, "content c", (*c.FieldB).FieldC)
	}
}

func Test_Parse_SinglePointerToStruct_Invalid(t *testing.T) {
	type substruct struct {
		FieldC int `key:"field_c"`
	}
	type configuration struct {
		FieldA string     `key:"field_a"`
		FieldB *substruct `key:"field_b"`
	}

	defer setEnvTemporarily("FIELD_A", "content a")()
	defer setEnvTemporarily("FIELD_B_FIELD_C", "ain't no integer here buddy")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrParseValue)
}

func Test_Parse_Struct_Invalid(t *testing.T) {
	type substruct struct {
		FieldC int `key:"field_c"`
	}
	type configuration struct {
		FieldA string    `key:"field_a"`
		FieldB substruct `key:"field_b"`
	}

	defer setEnvTemporarily("FIELD_A", "content a")()
	defer setEnvTemporarily("FIELD_B_FIELD_C", "ain't no integer here buddy")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrParseValue)
}

func Test_Parse_SingleNilPointerToStruct(t *testing.T) {
	type substruct struct {
		FieldC string `key:"field_c"`
	}
	type configuration struct {
		FieldA string     `key:"field_a"`
		FieldB *substruct `key:"field_b"`
	}

	defer setEnvTemporarily("FIELD_A", "content a")()
	defer setEnvTemporarily("FIELD_B_FIELD_C", "content c")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "content a", c.FieldA)
		assert.Equal(t, "content c", (*c.FieldB).FieldC)
	}
}

func Test_Parse_PointerOfDoomToStruct(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
		FieldB **************struct {
			FieldC string `key:"field_c"`
		} `key:"field_b"`
	}

	defer setEnvTemporarily("FIELD_A", "content a")()
	defer setEnvTemporarily("FIELD_B_FIELD_C", "content c")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "content a", c.FieldA)
		assert.Equal(t, "content c", (**************c.FieldB).FieldC)
	}
}

func Test_Parse_String_Whitespace(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "   ")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "   ", c.FieldA)
	}
}

func Test_Parse_Bool_Valid(t *testing.T) {
	type configuration struct {
		True       bool `key:"true"`
		False      bool `key:"false"`
		TrueUpper  bool `key:"true_upper"`
		FalseUpper bool `key:"false_upper"`
	}

	defer setEnvTemporarily("TRUE", "true")()
	defer setEnvTemporarily("FALSE", "false")()
	defer setEnvTemporarily("TRUE_UPPER", "TRUE")()
	defer setEnvTemporarily("FALSE_UPPER", "FALSE")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, true, c.True)
		assert.Equal(t, false, c.False)
		assert.Equal(t, true, c.TrueUpper)
		assert.Equal(t, false, c.FalseUpper)
	}
}

func Test_Parse_Bool_Invalid(t *testing.T) {
	type configuration struct {
		Bool bool `key:"bool"`
	}

	defer setEnvTemporarily("BOOL", "cheese")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrParseValue)
}

func Test_Parse_Complex64_Unsupported(t *testing.T) {
	type configuration struct {
		FieldA complex64 `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "value")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrUnsupportedFieldType)
}

func Test_Parse_Complex128_Unsupported(t *testing.T) {
	type configuration struct {
		FieldA complex128 `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "value")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrUnsupportedFieldType)
}

func Test_Parse_Int_Valid(t *testing.T) {
	type configuration struct {
		Min int `key:"min"`
		Max int `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", math.MinInt))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", math.MaxInt))()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, math.MinInt, c.Min)
		assert.Equal(t, math.MaxInt, c.Max)
	}
}

func Test_Parse_Int8_Valid(t *testing.T) {
	type configuration struct {
		Min int8 `key:"min"`
		Max int8 `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", math.MinInt8))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", math.MaxInt8))()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, math.MinInt8, c.Min)
		assert.Equal(t, math.MaxInt8, c.Max)
	}
}

func Test_Parse_Int16_Valid(t *testing.T) {
	type configuration struct {
		Min int16 `key:"min"`
		Max int16 `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", math.MinInt16))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", math.MaxInt16))()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, math.MinInt16, c.Min)
		assert.Equal(t, math.MaxInt16, c.Max)
	}
}

func Test_Parse_Int32_Valid(t *testing.T) {
	type configuration struct {
		Min int32 `key:"min"`
		Max int32 `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", math.MinInt32))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", math.MaxInt32))()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, math.MinInt32, c.Min)
		assert.Equal(t, math.MaxInt32, c.Max)
	}
}

func Test_Parse_Int64_Valid(t *testing.T) {
	type configuration struct {
		Min int64 `key:"min"`
		Max int64 `key:"max"`
	}

	defer setEnvTemporarily("MIN", fmt.Sprintf("%d", math.MinInt64))()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", math.MaxInt64))()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, math.MinInt64, c.Min)
		assert.Equal(t, math.MaxInt64, c.Max)
	}
}

func Test_Parse_Uint_Valid(t *testing.T) {
	type configuration struct {
		Min uint `key:"min"`
		Max uint `key:"max"`
	}

	defer setEnvTemporarily("MIN", "0")()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", uint(math.MaxUint)))()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, uint(0), c.Min)
		assert.Equal(t, uint(math.MaxUint), c.Max)
	}
}

func Test_Parse_Uint8_Valid(t *testing.T) {
	type configuration struct {
		Min uint8 `key:"min"`
		Max uint8 `key:"max"`
	}

	defer setEnvTemporarily("MIN", "0")()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", uint8(math.MaxUint8)))()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, uint8(0), c.Min)
		assert.Equal(t, uint8(math.MaxUint8), c.Max)
	}
}

func Test_Parse_Uint16_Valid(t *testing.T) {
	type configuration struct {
		Min uint16 `key:"min"`
		Max uint16 `key:"max"`
	}

	defer setEnvTemporarily("MIN", "0")()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", uint16(math.MaxUint16)))()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, uint16(0), c.Min)
		assert.Equal(t, uint16(math.MaxUint16), c.Max)
	}
}

func Test_Parse_Uint32_Valid(t *testing.T) {
	type configuration struct {
		Min uint32 `key:"min"`
		Max uint32 `key:"max"`
	}

	defer setEnvTemporarily("MIN", "0")()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", uint32(math.MaxUint32)))()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, uint32(0), c.Min)
		assert.Equal(t, uint32(math.MaxUint32), c.Max)
	}
}

func Test_Parse_Uint64_Valid(t *testing.T) {
	type configuration struct {
		Min uint64 `key:"min"`
		Max uint64 `key:"max"`
	}

	defer setEnvTemporarily("MIN", "0")()
	defer setEnvTemporarily("MAX", fmt.Sprintf("%d", uint64(math.MaxUint64)))()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, uint64(0), c.Min)
		assert.Equal(t, uint64(math.MaxUint64), c.Max)
	}
}

func Test_Parse_Float32_Valid(t *testing.T) {
	type configuration struct {
		Float float32 `key:"float"`
	}

	var floatValue float32 = 5.5
	bytes, _ := json.Marshal(floatValue)
	defer setEnvTemporarily("FLOAT", string(bytes))()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, floatValue, c.Float)
	}
}

func Test_Parse_Float64_Valid(t *testing.T) {
	type configuration struct {
		Float float64 `key:"float"`
	}

	var floatValue float64 = 5.5
	bytes, _ := json.Marshal(floatValue)
	defer setEnvTemporarily("FLOAT", string(bytes))()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, floatValue, c.Float)
	}
}

func Test_Parse_Float32_Invalid(t *testing.T) {
	type configuration struct {
		Float float32 `key:"float"`
	}

	defer setEnvTemporarily("FLOAT", "5.5no float here")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrParseValue)
}

func Test_Parse_Float64_Invalid(t *testing.T) {
	type configuration struct {
		Float float64 `key:"float"`
	}

	defer setEnvTemporarily("FLOAT", "5.5no float here")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrParseValue)
}

func Test_Parse_Int_Invalid(t *testing.T) {
	type configuration struct {
		FieldA int `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "10no int here")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrParseValue)
}

func Test_Parse_Uint_Invalid(t *testing.T) {
	type configuration struct {
		FieldA uint `key:"field_a"`
	}

	defer setEnvTemporarily("FIELD_A", "10no int here")()
	var c configuration
	err := yagcl.New[configuration]().Add(Source()).Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrParseValue)
}

func Test_Parse_DefaultValue_String(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a"`
	}

	c := configuration{
		FieldA: "i am the default",
	}
	err := yagcl.
		New[configuration]().
		Add(Source()).
		Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "i am the default", c.FieldA)
	}
}

func Test_Parse_DefaultValue_Int(t *testing.T) {
	type configuration struct {
		FieldA int `key:"field_a"`
	}

	c := configuration{
		FieldA: 1,
	}
	err := yagcl.
		New[configuration]().
		Add(Source()).
		Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, c.FieldA)
	}
}

func Test_Parse_DefaultValueZero_Int_Required(t *testing.T) {
	type configuration struct {
		FieldA int `key:"field_a" required:"true"`
	}

	c := configuration{
		FieldA: 0,
	}
	err := yagcl.
		New[configuration]().
		Add(Source()).
		Parse(&c)
	// If 0 is desired to be a valid value, a pointer should be used.
	// Alternatively remove "required:"true"".
	assert.ErrorIs(t, err, yagcl.ErrValueNotSet)
}

func Test_Parse_RequiredValue_Missing_String(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a" required:"true"`
	}

	var c configuration
	err := yagcl.
		New[configuration]().
		Add(Source()).
		Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrValueNotSet)
}

func Test_Parse_RequiredValue_EmptyDefault_String(t *testing.T) {
	type configuration struct {
		FieldA string `key:"field_a" required:"true"`
	}

	var c configuration
	err := yagcl.
		New[configuration]().
		Add(Source()).
		Parse(&c)
	assert.ErrorIs(t, err, yagcl.ErrValueNotSet)
}
