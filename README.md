# yagcl

[![Go Reference](https://pkg.go.dev/badge/github.com/Bios-Marcel/yagcl.svg)](https://pkg.go.dev/github.com/Bios-Marcel/yagcl)
[![Build and Tests](https://github.com/Bios-Marcel/yagcl/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/Bios-Marcel/yagcl/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/Bios-Marcel/yagcl/branch/master/graph/badge.svg?token=BPGE55G1AX)](https://codecov.io/gh/Bios-Marcel/yagcl)

This libraries aim is to provide a powerful and dynamic way to provide
configuration for your application.

## Why

The thing that other libraries were lacking is the ability to parse different
formats, allow merging them (for example override a setting via environment variables).
Additionally I wanna be able to specify validation, parsing, defaults and constraints
all in a central place: the field tags.

The aim is to support all standard datatypes and allow nested structs with specified
sub prefixes as well as one main prefix.

Additionally it is planned for the consumer of the library to be able to
validate a struct, essentially making sure it does't contain nonsensical
combinations of tags.

For example, the following wouldn't really make sense, since defining a key
for an ignored field has no effect and will therefore result in an error:

```go
type Configuration struct {
    Field string `key:"field" ignore:"true"`
}
```

If there's already a library that does ALL of this, feel free to tell me and I'll
delete the repository 😉.

## Modules

| Name                  | Repo                                               | Docs                                                                                                                                                 | Status  |
| --------------------- | -------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| Environment Variables | [GitHub](https://github.com/Bios-Marcel/yagcl-env) | [![Go Reference env package](https://pkg.go.dev/badge/github.com/Bios-Marcel/yagcl-env.svg)](https://pkg.go.dev/github.com/Bios-Marcel/yagcl-env)    | WIP     |
| JSON                  | [GitHub](https://github.com/Bios-Marcel/yagcl-env) | [![Go Reference json package](https://pkg.go.dev/badge/github.com/Bios-Marcel/yagcl-json.svg)](https://pkg.go.dev/github.com/Bios-Marcel/yagcl-json) | Planned |
| .env                  | -                                                  | -                                                                                                                                                    | Planned |

Also check out the [Roadmap](#roadmap) for more detailed information.

## Contribution

This library is separated into multiple modules. The main module and additional
modules for each supported source. This allows you to only specify certain
sources in your go-mod, keeping your dependency tree small. Additionally it
makes navigating the code base easier.

If you wish to contribute a new source, please create a corresponding
submodule.

## Examples

An example configuration usage may look something like this:

```go
import (
    yagcl_env "github.com/Bios-Marcel/yagcl-env"
    yagcl_json "github.com/Bios-Marcel/yagcl-json"
)

type Configuration struct {
    // The `key` here is used to define the JSON name for example. But the
    // environment variable names are the same, but uppercased.
    Host string `key:"host" required:"true"`
    Post int    `key:"port" required:"true"`
    // If you don't wish to export a field, you have to ignore it.
    // If it isn't ignored and doesn't have an explicit key, you'll
    // get an error, as this indicates a bug. The reason we don't
    // auto-generate a key is that this could result in unstable promises
    // as the variable name could change and break loading of old files.
    DontLoad    int               `ignore:"true"`
    // Nested structs are special, as they may not be part of your actual
    // configuration in case you are using environment variables, but will
    // be if you are using a JSON file. Either way, these also require the
    // key tag, as we are otherweise unable to build the names for its fields.
    KafkaServer KafkaServerConfig `key:"kafka"`
}

type KafkaServerConfig struct {
    //Alternatively you can define them explicitly. The same goes for json names.
    Host              string        `json:"host" env:"HOST" required:"true"`
    Port              int           `json:"port" env:"PORT" required:"true"`
    ConnectionTimeout time.Duration `json:"connection_timeout" env:"CONNECTION_TIMEOUT" required:"false"`
}

func LoadConfig() error {
    //Defaults should simply be defined on struct creation.
    configuration := Configuration{
        KafkaServer: KafkaServerConfig{
            Host:              "localhost",
            Port:              1234,
            ConnectionTimeout: time.Second * 10,
        },
    }
    err := yagcl.
        New[Configuration]()
        //This allows ordering when using override, so you can have something like this.
        Add(yagcl_json.Source("/etc/myapp/config.json").Must()).
        Add(yagcl_env.Source().Prefix("MY_APP_")).
        Add(yagcl_json.Source("~/.config/config.json")).
        AllowOverride().
        Parse(&configuration)
    return err
}
```

The configuration loaded by this would look like this:

```json
{
    "host": "localhost",
    "port": 1234,
    "kafka": {
        "host": "123.123.123.123",
        "port": 9092,
        "connection_timeout": "10s"
    }
}
```

Or this when loading environment variables:

```env
MY_APP_HOST=localhost
MY_APP_PORT=1234
MY_APP_KAFKA_HOST=123.123.123.123
MY_APP_KAFKA_PORT=9092
MY_APP_KAFKA_CONNECTION_TIMEOUT=10s
```

## Usage

**This library isn't stable / feature complete yet, even if it mostly works. The API might change any second ;)**

If you want to try it out anyway, simply `go get` the desired modules.

For example:

```shell
go get github.com/Bios-Marcel/yagcl-env
```

## Roadmap

- [x] Basic API
- [ ] General Features
  - [ ] Honor `required` tags
  - [ ] Validation of configuration struct
  - [ ] Functioning Override mechanism where a whole source is optional or only some fields
  > While overriding in general works, we'll error as soon as we are missing
  > one required value in any of the sources.
- [ ] Read JSON
  - [ ] Honor `key` tags
  - [x] Honor `json`tags
  - [ ] Honor `ignore` tags
- [ ] Read Environment variables
  - [x] Honor `key` tags
  - [x] Honor `env` tags
  - [x] Honor `ignore` tags
  - [ ] Type support
    - [x] int / uint
    - [x] float
    - [x] bool
    - [x] string
    - [x] struct
    - [x] pointer
    - [x] time.Duration
    - [ ] array
    - [ ] map
- [ ] Read .env files
  > Will share code with environment variables and should have the same progression.
