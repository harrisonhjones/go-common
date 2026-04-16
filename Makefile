.PHONY: test test-verbose cover lint fmt vet tidy clean

## Run all tests
test:
	go test ./...

## Run all tests with verbose output
test-verbose:
	go test -v ./...

## Run tests with coverage and open report
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## Run go vet
vet:
	go vet ./...

## Format all Go files
fmt:
	gofmt -w .

## Tidy module dependencies
tidy:
	go mod tidy

## Remove generated files
clean:
	rm -f coverage.out coverage.html
