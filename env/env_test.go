package env

import (
	"os"
	"testing"

	"github.com/Bios-Marcel/yagcl"
	"github.com/stretchr/testify/assert"
)

func Test_EventSource_InterfaceCompliance(t *testing.T) {
	var _ yagcl.Source = Source()
}

func Test_Parse_Env_Simple(t *testing.T) {
	type Configuration struct {
		FieldA string `key:"field_a"`
		FieldB string `env:"FIELD_B"`
	}

	os.Setenv("FIELD_A", "content a")
	os.Setenv("FIELD_B", "content b")
	var c Configuration
	err := yagcl.New[Configuration]().
		AddSource(Source().Prefix("TEST_")).
		Parse(&c)
	assert.NoError(t, err)
	assert.Equal(t, "content a", c.FieldA)
	assert.Equal(t, "content b", c.FieldB)
}
