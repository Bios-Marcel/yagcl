package env

type source struct {
	prefix string
}

// Source creates a source for environment variables of the current
// process.
func Source() *source {
	return &source{}
}

// Prefix specified the prefixes expected in environment variable keys.
// For example "PREFIX_FIELD_NAME".
func (s *source) Prefix(prefix string) *source {
	s.prefix = prefix
	return s
}

// KeyTag implements Source.Key.
func (s *source) KeyTag() string {
	return "env"
}

// Parse implements Source.Parse.
func (s *source) Parse(configurationStruct any) error {
	return nil
}
