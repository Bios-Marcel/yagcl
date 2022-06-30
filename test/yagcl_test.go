package test

import (
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
func (s *dummySource) Parse(configurationStruct any) (bool, error) {
	return true, nil
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
