package env

import (
	"encoding/json"
	"errors"
	"fmt"
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
	keyJoiner         func(string, string) string
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
		keyJoiner: func(s1, s2 string) string {
			if s1 == "" {
				return s2
			}
			if s2 == "" {
				return s1
			}
			// By default we want to use whatever keys we have, and join them
			// with underscores, preventing duplicate underscores.
			return strings.Trim(s1, "_") + "_" + strings.Trim(s2, "_")
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
// converted for this source. If you are setting this, you'll most liekly
// also have to set EnvSource.KeyJoiner(string,string) string.
// Note that calling this isn't required, as there's a best practise default
// behaviour.
func (s *EnvSource) KeyValueConverter(keyValueConverter func(string) string) *EnvSource {
	s.keyValueConverter = keyValueConverter
	return s
}

func (s *EnvSource) KeyJoiner(keyJoiner func(string, string) string) *EnvSource {
	s.keyJoiner = keyJoiner
	return s
}

// KeyTag implements Source.Key.
func (s *EnvSource) KeyTag() string {
	return "env"
}

// Parse implements Source.Parse.
func (s *EnvSource) Parse(configurationStruct any) error {
	return s.parse(s.prefix, reflect.Indirect(reflect.ValueOf(configurationStruct)))
}

func (s *EnvSource) parse(envPrefix string, structValue reflect.Value) error {
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
		joinedEnvKey := s.keyJoiner(envPrefix, envKey)
		envValue, set := os.LookupEnv(joinedEnvKey)

		//FIXME Do we need to differentiate here?
		if !set || envValue == "" {
			envValue, _ = structField.Tag.Lookup("default")
		}

		parsed, err := parseValue(structField.Name, structField.Type, envValue)
		if err != nil {
			if err == errEmbeddedStructDetected {
				if err := s.parse(joinedEnvKey, value); err != nil {
					return err
				}
				continue
			}
			return err
		}

		if parsed.IsZero() && strings.EqualFold(structField.Tag.Get("required"), "true") {
			return fmt.Errorf("environment variable '%s' not set correctly: %w", envKey, yagcl.ErrValueNotSet)
		}

		if value.Kind() == reflect.Pointer {
			//Create as many values as we have pointers pointing to things.
			var pointers []reflect.Value
			lastPointer := reflect.New(value.Type().Elem())
			pointers = append(pointers, lastPointer)
			for lastPointer.Elem().Kind() == reflect.Pointer {
				lastPointer = reflect.New(lastPointer.Elem().Type().Elem())
				pointers = append(pointers, lastPointer)
			}

			pointers[len(pointers)-1].Elem().Set(parsed)
			for i := len(pointers) - 2; i >= 0; i-- {
				pointers[i].Elem().Set(pointers[i+1])
			}
			value.Set(pointers[0])
		} else {
			value.Set(parsed)
		}
	}

	return nil
}

// errEmbeddedStructDetected is abused internally to detect that we need to
// recurse. This error should never reach the outer world.
var errEmbeddedStructDetected = errors.New("embedded struct detected")

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
			// Technically dead code right now, but we'll leave it in, as I am
			// unsure how the API will develop. Maybe overriding of keys should
			// be allowed to prevent clashing with other libraries?
			return "", fmt.Errorf("standard tag '%s' has not been set: %w", yagcl.DefaultKeyTagName, yagcl.ErrExportedFieldMissingKey)
		}
		envKey = s.keyValueConverter(envKey)
	}
	return envKey, nil
}

func parseValue(fieldName string, fieldType reflect.Type, envValue string) (reflect.Value, error) {
	switch fieldType.Kind() {
	case reflect.String:
		{
			return reflect.ValueOf(envValue), nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			value, errParse := strconv.ParseInt(envValue, 10, int(fieldType.Size())*8)
			if errParse != nil {
				return reflect.Value{}, fmt.Errorf("value '%s' isn't parsable as an '%s' for field '%s': %w", envValue, fieldType.String(), fieldName, yagcl.ErrParseValue)
			}
			return reflect.ValueOf(value).Convert(fieldType), nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			value, errParse := strconv.ParseUint(envValue, 10, int(fieldType.Size())*8)
			if errParse != nil {
				return reflect.Value{}, fmt.Errorf("value '%s' isn't parsable as an '%s' for field '%s': %w", envValue, fieldType.String(), fieldName, yagcl.ErrParseValue)
			}
			return reflect.ValueOf(value).Convert(fieldType), nil
		}
	case reflect.Float32, reflect.Float64:
		{
			// We use the stdlib json encoder here, since there seems to be
			// special behaviour.
			var value float64
			if errParse := json.Unmarshal([]byte(envValue), &value); errParse != nil {
				return reflect.Value{}, fmt.Errorf("value '%s' isn't parsable as an '%s' for field '%s': %w", envValue, fieldType.String(), fieldName, yagcl.ErrParseValue)
			}
			return reflect.ValueOf(value).Convert(fieldType), nil
		}
	case reflect.Bool:
		{
			boolValue := strings.EqualFold(envValue, "true")
			// FIXME Allow enabling lax-behaviour?
			// Instead of assuming everything != true equals false, we assume
			// that the value is unintentionally wrong and return an error.
			if !boolValue && !strings.EqualFold(envValue, "false") {
				return reflect.Value{}, fmt.Errorf("value '%s' isn't parsable as a '%s' for field '%s': %w", envValue, fieldType.String(), fieldName, yagcl.ErrParseValue)
			}
			return reflect.ValueOf(boolValue), nil
		}
	case reflect.Struct:
		{
			return reflect.Value{}, errEmbeddedStructDetected
		}
	case reflect.Pointer:
		{
			return parseValue(fieldName, extractNonPointerFieldType(fieldType), envValue)
		}
	case reflect.Complex64, reflect.Complex128:
		{
			// Complex isn't supported, as for example it also isn't supported
			// by the stdlib json encoder / decoder.
			return reflect.Value{}, fmt.Errorf("type '%s' isn't supported and won't ever be: %w", fieldName, yagcl.ErrUnsupportedFieldType)
		}
	default:
		{
			return reflect.Value{}, fmt.Errorf("type '%s': %w", fieldName, yagcl.ErrUnsupportedFieldType)
		}
	}
}

func extractNonPointerFieldType(fieldType reflect.Type) reflect.Type {
	if fieldType.Kind() != reflect.Pointer {
		return fieldType
	}

	return extractNonPointerFieldType(fieldType.Elem())
}
