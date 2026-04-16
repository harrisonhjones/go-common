# go-common

Common helpers / utils for building go libraries and applications.

## Packages

### `bind`

Bind external data sources to Go structs using struct tags.

#### `bind.FromEnv`

Populates struct fields from environment variables.

```go
type Config struct {
    Host string `env:"APP_HOST" envDefault:"localhost"`
    Port int    `env:"APP_PORT" envDefault:"8080"`
}

var cfg Config
if err := bind.FromEnv(context.Background(), &cfg); err != nil {
    log.Fatal(err)
}
```

Supported tags:

- `env:"VAR_NAME"` — environment variable to read
- `envDefault:"value"` — fallback if the variable is unset or empty
- `envFromNoOverwrite:"true"` — only bind if the field is zero-valued

#### `bind.FromFlags`

Populates struct fields from command line flags.

```go
type Config struct {
    Host    string `flag:"host" flagDefault:"localhost" flagUsage:"server host"`
    Verbose bool   `flag:"verbose"`
}

var cfg Config
if err := bind.FromFlags(context.Background(), os.Args[1:], &cfg); err != nil {
    log.Fatal(err)
}
```

Supported tags:

- `flag:"name"` — flag name
- `flagDefault:"value"` — default value if flag is not provided
- `flagUsage:"text"` — usage description for help output
- `flagFromNoOverwrite:"true"` — only bind if the field is zero-valued

#### Supported types

`string`, `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`,
`uint16`, `uint32`, `uint64`, `float32`, `float64`, `[]string`

#### Normalizer / Validator

Structs can implement `Normalize(context.Context)` and/or
`Validate(context.Context) error` to run post-bind logic:

```go
func (c *Config) Normalize(ctx context.Context) {
    c.Host = strings.ToLower(c.Host)
}

func (c *Config) Validate(ctx context.Context) error {
    if c.Port == 0 {
        return fmt.Errorf("port is required")
    }
    return nil
}
```

## Development

```
make test           # run tests
make test-verbose   # run tests with verbose output
make cover          # generate coverage report
make vet            # run go vet
make fmt            # format code
make tidy           # tidy module deps
make clean          # remove generated files
```
