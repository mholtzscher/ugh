# ugh justfile
# Run `just` to see available recipes

# Default recipe: list available commands
default:
    @just --list

# Build the binary (development)
build:
    go build

# Build the binary (release, stripped)
build-release:
    go build -ldflags="-s -w"

# Run the binary with arguments (use `just -- run <args>` for flags)
run *ARGS:
    go run . {{ARGS}}

# Run all tests
test:
    go test -p 1 ./...

# Run tests with verbose output
test-verbose:
    go test -v ./...

# Run a specific test by name
test-run NAME:
    go test -run {{NAME}} ./...

# Run tests for a specific package
test-pkg PKG:
    go test ./{{PKG}}

# Format code
fmt:
    go fmt ./...

# Run static analysis (excluding generated ANTLR code)
vet:
    go vet $(go list ./... | grep -v /nlp/antlr/parser)

# Run comprehensive linting
lint:
    golangci-lint run

# Generate code (sqlc + antlr)
generate:
   sqlc generate
   @command -v antlr4 >/dev/null 2>&1 && (cd internal/nlp/antlr && antlr4 -Dlanguage=Go -package parser -visitor -listener -o parser UghLexer.g4 UghParser.g4) || echo "antlr4 not found, skipping ANTLR generation (add antlr to devShell or run: nix shell nixpkgs#antlr -c just generate)"

# Run all checks (generate, format, vet, lint, test, tidy, gomod2nix)
check: generate fmt vet lint test tidy gomod2nix

# Update Go dependencies
update-deps:
    go get -u ./...
    go mod tidy
    gomod2nix > gomod2nix.toml

# Tidy go modules
tidy:
    go mod tidy

# Update gomod2nix.toml after dependency changes
gomod2nix:
    gomod2nix > gomod2nix.toml

# Run govulncheck security scanner
govulncheck:
    govulncheck ./...

# Validate template consistency
cruft-check:
    cruft check

# Show template differences
cruft-diff:
    cruft diff

# Update to latest template
cruft-update:
    cruft update

# Run with nix
nix-run *ARGS:
    nix run -- {{ARGS}}
