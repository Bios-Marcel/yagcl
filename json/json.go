package json

import (
	"encoding/json"
	"os"

	"github.com/Bios-Marcel/yagcl"
)

type source struct {
	must bool
	path string
}

// Source creates a source for a JSON file.
func Source(path string) *source {
	return &source{path: path}
}

// Must indicates that this Source will return an error during parsing if no
// parsable data can be found.
func (s *source) Must() *source {
	s.must = true
	return s
}

// KeyTag implements Source.Key.
func (s *source) KeyTag() string {
	return "json"
}

// Parse implements Source.Parse.
func (s *source) Parse(configurationStruct any) error {
	file, errOpen := os.OpenFile(s.path, os.O_RDONLY, os.ModePerm)
	if os.IsNotExist(errOpen) {
		if s.must {
			return yagcl.ErrSourceNotFound
		}
		return nil
	}

	return json.
		NewDecoder(file).
		Decode(configurationStruct)
}
