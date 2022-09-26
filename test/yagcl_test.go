package test

import (
	"reflect"
	"testing"

	"github.com/Bios-Marcel/yagcl"
	"github.com/stretchr/testify/assert"
)

type dummySource struct {
}

// KeyTag implements Source.Key.
func (s *dummySource) KeyTag() string {
	return "dummy"
}

// Parse implements Source.Parse.
func (s *dummySource) Parse(parsingCompanion yagcl.ParsingCompanion, configurationStruct any) (bool, error) {
	structValue := reflect.Indirect(reflect.ValueOf(configurationStruct))
	structType := structValue.Type()
	for i := 0; i < structValue.NumField(); i++ {
		structField := structType.Field(i)
		if !parsingCompanion.IncludeField(structField) {
			continue
		}

		value := structValue.Field(i)
		value.Set(reflect.ValueOf(parsingCompanion.ExtractFieldKey(structField)))
	}
	return true, nil
}

func Test_Parse_InferFieldKey(t *testing.T) {
	type config struct {
		Field string
	}
	var c config
	err := yagcl.
		New[config]().
		Add(&dummySource{}).
		InferFieldKeys().
		Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "field", c.Field)
	}

	c = config{}
	err = yagcl.
		New[config]().
		Add(&dummySource{}).
		Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "", c.Field)
	}
}

func Test_Parse_AdditionalKeyTags(t *testing.T) {
	type config struct {
		Field string `kek:"oof"`
	}
	var c config
	err := yagcl.
		New[config]().
		Add(&dummySource{}).
		AdditionalKeyTags("kek").
		Parse(&c)
	if assert.NoError(t, err) {
		assert.Equal(t, "oof", c.Field)
	}
}

func Test_Parse_PassNilPointer(t *testing.T) {
	type config struct{}
	err := yagcl.
		New[config]().
		Add(&dummySource{}).
		Parse(nil)
	assert.ErrorIs(t, err, yagcl.ErrInvalidConfiguraionPointer)
}

func Test_Parse_PassNilPointerVariable(t *testing.T) {
	type config struct{}
	var nilCfg *config
	err := yagcl.
		New[config]().
		Add(&dummySource{}).
		Parse(nilCfg)
	assert.ErrorIs(t, err, yagcl.ErrInvalidConfiguraionPointer)
}

func Test_Parse_PassZeroStruct(t *testing.T) {
	type config struct{}
	cfg := &config{}
	err := yagcl.
		New[config]().
		Add(&dummySource{}).
		Parse(cfg)
	if assert.NoError(t, err) {
		assert.NotNil(t, cfg, "We expect initialise any pointer to a struct thats nil")
	}
}
