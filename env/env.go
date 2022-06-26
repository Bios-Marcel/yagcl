package env

import (
	"fmt"
	"math/bits"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/Bios-Marcel/yagcl"
)

// DO NOT CREATE INSTANCES MANUALLY, THIS IS ONLY PUBLIC IN ORDER FOR GODOC
// TO RENDER AVAILABLE FUNCTIONS.
type EnvSource struct {
	prefix            string
	keyValueConverter func(string) string
}

// Source creates a source for environment variables of the current
// process.
func Source() *EnvSource {
	return &EnvSource{
		keyValueConverter: func(s string) string {
			// Since by default we expect keys to be of
			// format `word_word_...`, we just uppercase everything to meet
			// the defacto standard of environment variables.
			return strings.ToUpper(s)
		},
	}
}

// Prefix specified the prefixes expected in environment variable keys.
// For example "PREFIX_FIELD_NAME".
func (s *EnvSource) Prefix(prefix string) *EnvSource {
	s.prefix = prefix
	return s
}

// KeyValueConverter defines how the yagcl.DefaultKeyTagName value should be
// converted for this source. Note that calling this isn't required, as there's
// a best practise default behaviour.
func (s *EnvSource) KeyValueConverter(keyValueConverter func(string) string) *EnvSource {
	s.keyValueConverter = keyValueConverter
	return s
}

// KeyTag implements Source.Key.
func (s *EnvSource) KeyTag() string {
	return "env"
}

// Parse implements Source.Parse.
func (s *EnvSource) Parse(configurationStruct any) error {
	structValue := reflect.Indirect(reflect.ValueOf(configurationStruct))
	structType := structValue.Type()
	for i := 0; i < structValue.NumField(); i++ {
		structField := structType.Field(i)
		// By default, all exported fiels are not ignored and all exported
		// fields are. Unexported fields can't be un-ignored though.
		if strings.EqualFold(structField.Tag.Get("ignore"), "true") {
			continue
		}

		value := structValue.Field(i)
		envKey, err := s.extractEnvKey(value, structField)
		if err != nil {
			return err
		}
		envValue, set := os.LookupEnv(envKey)

		//FIXME Do we need to differentiate here?
		if !set || envValue == "" {
			envValue, _ = structField.Tag.Lookup("default")
		}

		parsed, err := parseValue(structField, envValue)
		if err != nil {
			return err
		}

		if parsed.IsZero() && strings.EqualFold(structField.Tag.Get("required"), "true") {
			return fmt.Errorf("environment variable '%s' not set correctly: %w", envKey, yagcl.ErrValueNotSet)
		}

		value.Set(*parsed)
	}

	return nil
}

func (s *EnvSource) extractEnvKey(value reflect.Value, structField reflect.StructField) (string, error) {
	var (
		envKey string
		tagSet bool
	)
	customKeyTag := s.KeyTag()
	if customKeyTag != "" {
		envKey, tagSet = structField.Tag.Lookup(customKeyTag)
	}
	if !tagSet {
		envKey, tagSet = structField.Tag.Lookup(yagcl.DefaultKeyTagName)
		if !tagSet {
			if customKeyTag != "" {
				return "", fmt.Errorf("neither tag '%s' nor the standard tag '%s' have been set: %w", customKeyTag, yagcl.DefaultKeyTagName, yagcl.ErrExportedFieldMissingKey)
			}
			return "", fmt.Errorf("standard tag '%s' has not been set: %w", yagcl.DefaultKeyTagName, yagcl.ErrExportedFieldMissingKey)
		}
		envKey = s.keyValueConverter(envKey)
	}
	return s.prefix + envKey, nil
}

func parseValue(structField reflect.StructField, envValue string) (*reflect.Value, error) {
	var parsed reflect.Value
	switch structField.Type.Kind() {
	case reflect.String:
		{
			parsed = reflect.ValueOf(envValue)
		}
	case reflect.Int:
		{
			value, errParse := strconv.ParseInt(envValue, 10, bits.UintSize)
			if errParse != nil {
				return nil, fmt.Errorf("value '%s' isn't parsable as an '%s' for field '%s': %w", envValue, reflect.Int.String(), structField.Name, yagcl.ErrParseValue)
			}
			parsed = reflect.ValueOf(int(value))
		}
	default:
		{
			return nil, fmt.Errorf("type '%s': %w", structField.Name, yagcl.ErrUnsupportedFieldType)
		}
	}

	return &parsed, nil
}
