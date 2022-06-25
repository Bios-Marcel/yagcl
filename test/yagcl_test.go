package yagcl

import (
	"testing"

	"github.com/Bios-Marcel/yagcl"
	"github.com/stretchr/testify/assert"
)

func Test_EventSource_InterfaceCompliance(t *testing.T) {
	var _ yagcl.Source = yagcl.EventSource()
}

func Test_JSONSource_InterfaceCompliance(t *testing.T) {
	var _ yagcl.Source = yagcl.JSONSource("")
}

func Test_Parse_Simple(t *testing.T) {
	type Configuration struct {
		Field string `json:"field"`
	}
	var c Configuration
	err := yagcl.New[Configuration]().
		AddSource(yagcl.JSONSource("./test.json").Must()).
		AddSource(yagcl.EventSource().Prefix("TEST")).
		Parse(&c)
	assert.NoError(t, err)
	assert.Equal(t, "content", c.Field)
}
