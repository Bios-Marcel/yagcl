package yagcl

import (
	"encoding/json"
	"errors"
	"os"
)

var (
	// ErrExpectAtLeastOneSource indicates that no configuration can be loaded
	// if there's not at least one source.
	ErrExpectAtLeastOneSource = errors.New("please define at least one source")
	// ErrSourceNotFound implies that a source has been added, but for example
	// the file it is supposed to read can't be found.
	ErrSourceNotFound = errors.New("the source could not find its target")
)

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

type envSource struct {
	prefix string
}

// EventSource creates a source for environment variables of the current
// process.
func EventSource() *envSource {
	return &envSource{}
}

// Prefix specified the prefixes expected in environment variable keys.
// For example "PREFIX_FIELD_NAME".
func (s *envSource) Prefix(prefix string) *envSource {
	s.prefix = prefix
	return s
}

// KeyTag implements Source.Key.
func (s *envSource) KeyTag() string {
	return "env"
}

// Parse implements Source.Parse.
func (s *envSource) Parse(configurationStruct any) error {
	return nil
}

type jsonSource struct {
	must bool
	path string
}

// JSONSource creates a source for a JSON file.
func JSONSource(path string) *jsonSource {
	return &jsonSource{path: path}
}

// Must indicates that this Source will return an error during parsing if no
// parsable data can be found.
func (s *jsonSource) Must() *jsonSource {
	s.must = true
	return s
}

// KeyTag implements Source.Key.
func (s *jsonSource) KeyTag() string {
	return "json"
}

// Parse implements Source.Parse.
func (s *jsonSource) Parse(configurationStruct any) error {
	file, errOpen := os.OpenFile(s.path, os.O_RDONLY, os.ModePerm)
	if os.IsNotExist(errOpen) {
		if s.must {
			return ErrSourceNotFound
		}
		return nil
	}

	return json.
		NewDecoder(file).
		Decode(configurationStruct)
}

// AddSource adds a single source to read configuration from. This method can
// be called multiple times, adding multiple ordered sources. Whatever is
// added first is preferred. If AllowOverride() is called, all source will be
// parsed in the defined order.
func (y *YAGCL[T]) AddSource(source Source) *YAGCL[T] {
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

	// Build cached data required for parsing
	// for _, source := range y.sources {
	// }

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
