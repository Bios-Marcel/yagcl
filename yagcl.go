package yagcl

import (
	"errors"
	"reflect"
	"strings"
)

var (
	// ErrExpectAtLeastOneSource indicates that no configuration can be loaded
	// if there's not at least one source.
	ErrExpectAtLeastOneSource = errors.New("please define at least one source")
	// ErrSourceNotFound implies that a source has been added, but for example
	// the file it is supposed to read can't be found.
	ErrSourceNotFound = errors.New("the source could not find its target")
	// ErrInvalidConfiguraionPointer TODO
	ErrInvalidConfiguraionPointer = errors.New("invalid configuration pointer passed")

	// ErrExportedFieldMissingKey implies that a field is exported, but
	// doesn't define the key required for parsing it from a source. This
	// error can be omitted by setting `ignore:"true"`
	ErrExportedFieldMissingKey = errors.New("exported field is missing key definition")
	// ErrValueNotSet implies that no value could be found for a field, even
	// though it is required. The reason is either that there's no value, no
	// default value or a user error has occurred (e.g. a typo in the key).
	ErrValueNotSet = errors.New("required not set and no non-zero default value was found")
	// ErrParseValue implies that the value we attempted to parse is not in
	// the correct format to be assigned to its corresponding field. This is
	// most likely a user error.
	ErrParseValue = errors.New("value not parsable as type specified by field")
	// ErrUnsupportedFieldType implies that this library does NOT YET support
	// a certain field type.
	ErrUnsupportedFieldType = errors.New("unsupported field type")
)

// DefaultKeyTagName is the go annotation key name for specifying the default
// fieldname. The value is to be expected to use only lowercase letters and
// underscores, as this is deemed the best default for readability. Sources
// can then convert this to kebapCase, UPPER_CASE or whatever else the put
// their faith in.
const DefaultKeyTagName = "key"

// Source represents a source for configuration values. This might be backed
// by files of different formats, network access, environment variables or
// even the Windows registry. While YAGCL offers default sources, you can
// easily integrate your own source.
type Source interface {
	// Parse attempts retrieving data from the source and parsing it. This
	// should return true if anything was loaded. If nothing was loaded, but we
	// expected to load data, we should return an error. A realistic scenario
	// here would be that a Source is required instead of "load if found".
	Parse(ParsingCompanion, any) (bool, error)
	// KeyTag defines the golang struct field tag that defines a specific tag
	// for this source. This allows having one generic key, but overrides for
	// specific sources.
	KeyTag() string
}

type yagclImpl[T any] struct {
	sources        []Source
	allowOverride  bool
	inferFieldKeys bool
	keyTags        []KeyTag
}

// KeyTag describes a supporte that contains a field key. Additionally you
// can define a custom parser, that manipulates the content before using it.
// This can for example be useful if you want ot use a pre-existing tag, such
// as the standard libraries JSON tag, which may contain additional values
// separated by a comma, such as "omitempty".
type KeyTag struct {
	Name   string
	Parser func(string) string
}

// YAGCL defines the setup interface for YAGCL.
type YAGCL[T any] interface {
	// Add adds a single source to read configuration from. This method can
	// be called multiple times, adding multiple ordered sources. Whatever is
	// added first is preferred. If AllowOverride() is called, all source will be
	// parsed in the defined order.
	Add(source Source) YAGCL[T]
	// AllowOverride allows YAGCL to read from multiple sources. For example more
	// than one JSON file or a JSON file and the environment.
	AllowOverride() YAGCL[T]
	// InferFieldKeys activates automtic generation of a field key if none has
	// been defined by the programmers.
	InferFieldKeys() YAGCL[T]
	// AdditionalKeyTags defines tags other than `key`, which will then be used
	// by ParsingCompanion.ExtractFieldKey.
	AdditionalKeyTags(tags ...KeyTag) YAGCL[T]

	// Parse expects a pointer to a struct, which it'll attempt loading the
	// configuration into. Note that you'll first have to specify any type
	// type of configuration to be loaded.
	Parse(configurationStruct *T) error
}

// New creates a fresh instance of YAGCL. This is an alternative to using the
// global instance accessible via Global().
func New[T any]() YAGCL[T] {
	return &yagclImpl[T]{
		keyTags: []KeyTag{{Name: DefaultKeyTagName}},
	}
}

func (y *yagclImpl[T]) Add(source Source) YAGCL[T] {
	y.sources = append(y.sources, source)
	return y
}

func (y *yagclImpl[T]) AllowOverride() YAGCL[T] {
	y.allowOverride = true
	return y
}

func (y *yagclImpl[T]) InferFieldKeys() YAGCL[T] {
	y.inferFieldKeys = true
	return y
}

func (y *yagclImpl[T]) AdditionalKeyTags(tags ...KeyTag) YAGCL[T] {
	y.keyTags = append(y.keyTags, tags...)
	return y
}

func (y *yagclImpl[T]) Parse(configurationStruct *T) error {
	// While no sources would technically be fine, it doesn't really make
	// sense to call Parse at all then. Unless sources are somehow defined
	// dynamically, but that would go against the idea of this congfiguration
	// library.
	if len(y.sources) == 0 {
		return ErrExpectAtLeastOneSource
	}

	// Since a nil pointer can't be set to a zero-struct from inside a
	// function, we have to error here.
	if configurationStruct == nil {
		return ErrInvalidConfiguraionPointer
	}

	for _, source := range y.sources {
		loaded, err := source.Parse(y, configurationStruct)
		if err != nil {
			return err
		}

		if loaded && !y.allowOverride {
			break
		}
	}

	return nil
}

type ParsingCompanion interface {
	// IncludeField determines whether the field should be included in parsing.
	// This defines the standard for YAGCL and should be used by all sources.
	// A source may support additional rules that may even overwrite this ruleset.
	IncludeField(reflect.StructField) bool
	// ExtractFieldKey defines which identifier should be used for the given
	// field. This defines the standard identifier defined by YAGCL if no
	// specific identifier has been found for your source.
	ExtractFieldKey(reflect.StructField) string
}

func (y *yagclImpl[T]) IncludeField(structField reflect.StructField) bool {
	return structField.IsExported() && !strings.EqualFold(structField.Tag.Get("ignore"), "true")
}

func (y *yagclImpl[T]) ExtractFieldKey(structField reflect.StructField) string {
	for _, keyTag := range y.keyTags {
		tagName, isSet := structField.Tag.Lookup(keyTag.Name)
		if isSet && tagName != "" {
			if keyTag.Parser != nil {
				return keyTag.Parser(tagName)
			}
			return tagName
		}
	}

	if y.inferFieldKeys {
		return strings.ToLower(structField.Name)
	}

	return ""
}
