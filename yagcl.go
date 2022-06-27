package yagcl

import (
	"errors"
)

var (
	// ErrExpectAtLeastOneSource indicates that no configuration can be loaded
	// if there's not at least one source.
	ErrExpectAtLeastOneSource = errors.New("please define at least one source")
	// ErrSourceNotFound implies that a source has been added, but for example
	// the file it is supposed to read can't be found.
	ErrSourceNotFound = errors.New("the source could not find its target")

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

type YAGCL[T any] struct {
	sources       []Source
	allowOverride bool
}

// New creates a fresh instance of YAGCL. This is an alternative to using the
// global instance accessible via Global().
func New[T any]() *YAGCL[T] {
	return &YAGCL[T]{}
}

// Source represents a source for configuration values. This might be backed
// by files of different formats, network access, environment variables or
// even the windows registry. While YAGCL offers default sources, you can
// easily integrate your own source.
type Source interface {
	// Parse attempts retrieving data from the source and parsing it.
	Parse(any) error
	// KeyTag defines the golang struct field tag that defines a specific tag
	// for this source. This allows having one generic key, but overrides for
	// specific sources.
	KeyTag() string
}

// Add adds a single source to read configuration from. This method can
// be called multiple times, adding multiple ordered sources. Whatever is
// added first is preferred. If AllowOverride() is called, all source will be
// parsed in the defined order.
func (y *YAGCL[T]) Add(source Source) *YAGCL[T] {
	y.sources = append(y.sources, source)
	return y
}

// AllowOverride allows YAGCL to read from multiple sources. For example more
// than one JSON file or a JSON file and the environment.
func (y *YAGCL[T]) AllowOverride() *YAGCL[T] {
	y.allowOverride = true
	return y
}

// Parse expects a pointer to a struct, which it'll attempt loading the
// configuration into. Note that you'll first have to specify any type
// type of configuration to be loaded.
func (y *YAGCL[T]) Parse(configurationStruct *T) error {
	if len(y.sources) == 0 {
		return ErrExpectAtLeastOneSource
	}

	// Do actual parsing.
	for _, source := range y.sources {
		if err := source.Parse(configurationStruct); err != nil {
			return err
		}

		if !y.allowOverride {
			break
		}
	}

	return nil
}
