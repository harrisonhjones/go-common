// Package bind populates Go structs from external data sources using struct
// tags. It currently supports environment variables and command-line flags.
//
// # Environment Variables
//
// Use [FromEnv] to populate a struct from the process environment:
//
//	type Config struct {
//	    Host string `env:"APP_HOST" envDefault:"localhost"`
//	    Port int    `env:"APP_PORT" envDefault:"8080"`
//	}
//
//	var cfg Config
//	if err := bind.FromEnv(context.Background(), &cfg); err != nil {
//	    log.Fatal(err)
//	}
//
// Supported struct tags for [FromEnv]:
//
//   - env:"VAR_NAME" — environment variable to read.
//   - envDefault:"value" — fallback when the variable is unset or empty.
//   - envFromNoOverwrite:"true" — skip binding if the field already has a
//     non-zero value.
//
// # Command-Line Flags
//
// Use [FromFlags] to populate a struct from command-line arguments:
//
//	type Config struct {
//	    Host    string `flag:"host" flagDefault:"localhost" flagUsage:"server host"`
//	    Verbose bool   `flag:"verbose"`
//	}
//
//	var cfg Config
//	if err := bind.FromFlags(context.Background(), os.Args[1:], &cfg); err != nil {
//	    log.Fatal(err)
//	}
//
// Supported struct tags for [FromFlags]:
//
//   - flag:"name" — flag name on the command line.
//   - flagDefault:"value" — default value when the flag is not provided.
//   - flagUsage:"text" — usage description shown in help output.
//   - flagFromNoOverwrite:"true" — skip binding if the field already has a
//     non-zero value.
//
// # Supported Field Types
//
// Both [FromEnv] and [FromFlags] support the following field types:
//
//	string, bool,
//	int, int8, int16, int32, int64,
//	uint, uint8, uint16, uint32, uint64,
//	float32, float64,
//	[]string (comma-separated)
//
// # Post-Bind Hooks
//
// If the target struct implements [Normalizer], its Normalize method is called
// after all fields are set. If it implements [Validator], Validate is called
// after normalization and any returned error is propagated:
//
//	func (c *Config) Normalize(ctx context.Context) {
//	    c.Host = strings.ToLower(c.Host)
//	}
//
//	func (c *Config) Validate(ctx context.Context) error {
//	    if c.Port == 0 {
//	        return fmt.Errorf("port is required")
//	    }
//	    return nil
//	}
package bind
