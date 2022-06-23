# yagcl

This libraries aim is to provide a powerful and dynamic way to provide
configuration for your application.

The thing that other libraries were lacking is the ability to parse different
formats, allow merging them (for example override a setting via environment variables).
Additionally I wanna be able to specify validation, parsing, defaults and constraints
all in a central place, the field tags.

The aim is to support all standard datatypes and allow nested structs with specified
sub prefixes as well as one main prefix.

An example struct may look something like this:

```go
type Configuration struct {
	Host        string `json:"host" env:"HOST" required:"true"`
	Post        int    `json:"port" env:"PORT" required:"true"`
	KafkaServer struct {
		Host              string `json:"host" env:"HOST" default:"localhost" required:"true"`
		Port              int    `json:"port" env:"PORT" default:"1234" required:"true"`
		ConnectionTimeout int    `json:"connection_timeout" env:"CONNECTION_TIMEOUT" default:"10s" required:"false"`
	} `json:"kafka" env_prefix:"KAFKA_"`
}

func LoadConfig() error {
	var configuration Configuration
	err := yagcl.
		ParseJSON("config.json").
		ParseEnv().
		EnvPrefix("MYAPP").
		AllowOverride(true).
		Parse(&configuration)
	return err
}
```

If there's already a library that does ALL of this, feel free to tell me and I'll
delete the repository 😉.

## Usage

**DON'T, there's no code yet. Even if ther was, wait til there's a tagged and tested version.**

