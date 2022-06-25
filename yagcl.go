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
