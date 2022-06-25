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
	//The environment variable names are the same, but uppercased.
	Host        string `json:"host" required:"true"`
	Post        int    `json:"port" required:"true"`
	KafkaServer struct {
		//Alternatively you can define them explicitly.
		Host              string        `key:"host" default:"localhost" required:"true"`
		Host              string        `json:"host" env:"HOST" default:"localhost" required:"true"`
		Port              int           `json:"port" env:"PORT" default:"1234" required:"true"`
		ConnectionTimeout time.Duration `json:"connection_timeout" env:"CONNECTION_TIMEOUT" default:"10s" required:"false"`
		//Nested structs are an exception, as we need a prefix for each
		//struct to prevent clashing. If no prefix has been defined, it'll
		//be inferred from the fieldname.
	} `json:"kafka" env_prefix:"KAFKA_"`
}

func LoadConfig() error {
	var configuration Configuration
	err := yagcl.
		//This allows ordering when using override, so you can have something like this.
		AddSource(json.Source("/etc/myapp/config.json").Must()).
		AddSource(env.Source().Prefix("MY_APP_")).
		AddSource(json.Source("~/.config/config.json")).
		AllowOverride().
		Parse(&configuration)
	return err
}
```

The configuration loaded by this could look something like this:

```json
{
	"host": "localhost",
	"port": 1234,
	"kafka": {
		"host": "123.123.123.123",
		"port": 9092,
		"connectionTimeout": "10s"
	}
}
```

```env
MY_APP_HOST=localhost
MY_APP_PORT=1234
MY_APP_KAFKA_HOST=123.123.123.123
MY_APP_KAFKA_PORT=9092
MY_APP_KAFKA_CONNECTION_TIMEOUT=10s
```

If there's already a library that does ALL of this, feel free to tell me and I'll
delete the repository 😉.

## Usage

**DON'T, there's no code yet. Even if ther was, wait til there's a tagged and tested version.**

