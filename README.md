# go-common

A small collection of general-purpose helpers for building Go libraries and applications. The goal is to reduce boilerplate for patterns that come up repeatedly — pointer wrappers, must-or-panic, interactive prompts, and struct binding from external sources.

[![Go Reference](https://pkg.go.dev/badge/harrisonhjones.com/go-common.svg)](https://pkg.go.dev/harrisonhjones.com/go-common)
[![Go Report Card](https://goreportcard.com/badge/harrisonhjones.com/go-common)](https://goreportcard.com/report/harrisonhjones.com/go-common)

## Packages

- `harrisonhjones.com/go-common` — pointer helpers, `Must`, and terminal prompts
- `harrisonhjones.com/go-common/bind` — bind env vars and CLI flags to structs via struct tags

Full API documentation is on [pkg.go.dev](https://pkg.go.dev/harrisonhjones.com/go-common).

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
