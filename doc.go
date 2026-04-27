// Package common provides general-purpose helper functions for building
// Go libraries and applications.
//
// # Pointer Helpers
//
// [Pointer] and [Value] simplify working with pointer types, which is
// especially useful when constructing struct literals for APIs that require
// pointer fields:
//
//	name := common.Pointer("alice")  // *string pointing to "alice"
//	val  := common.Value(name)       // "alice"
//	val   = common.Value[string](nil) // "" (zero value of string)
//
// # Must
//
// [Must] unwraps a (T, error) pair, panicking with the error if it is non-nil.
// It is intended for initialization code or tests where an error is unrecoverable:
//
//	tpl := common.Must(template.New("t").Parse("Hello, {{.}}!"))
//
// # Prompt
//
// [Prompter] provides simple interactive terminal prompts, reading from an
// [io.Reader] and writing to an [io.Writer]. [DefaultPrompter] uses stdin/stdout.
// Package-level [Prompt], [PromptWithDefault], and [PromptContinue] delegate
// to [DefaultPrompter]:
//
//	name, err := common.Prompt("Your name")
//	name, err  = common.PromptWithDefault("Your name", "World") // returns "World" if empty
//	ok, err   := common.PromptContinue(true)                    // defaults to yes
//
// For testing, construct a [Prompter] with custom reader/writer via [NewPrompter].
//
// # Sub-packages
//
//   - [harrisonhjones.com/go-common/bind] — bind environment variables and
//     command-line flags to Go structs using struct tags.
package common
