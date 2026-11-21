# Quick start

Use the Go version from go.mod (currently Go 1.27.0 minimum). From this checkout:

```sh
make setup
make format
make check-format
make check
```

## Check or rewrite explicit files

After v0.1.0 is published, install with
`go install github.com/longxi-go/golayout/cmd/golayout@v0.1.0`, or use
`go get -tool github.com/longxi-go/golayout/cmd/golayout@v0.1.0` in a tool module
and invoke `go tool golayout`. From this checkout:

```sh
go run ./cmd/golayout -lang go1.27 path/to/file.go
go run ./cmd/golayout -w -lang go1.27 path/to/file.go
```

The first command is read-only: layout differences return a nonzero exit status.
`-w` applies the shared policy. Commands require at least one explicit file;
repository discovery remains the caller's responsibility.

## Integrate a generator

Import `github.com/longxi-go/golayout/pkg/format`:

```go
formatted, err := format.Source(source, format.Options{
    LangVersion: "go1.27",
})
```

This is a fragment: `source` is a byte slice owned by the generator. Handle the
error before validating or writing `formatted`. Source input bytes are not mutated.
Derive the target language from the consuming module's go.mod; the literal above
only illustrates the option. Do not infer it from the generator's own toolchain.

Use the public API before type validation/output. Consumers should delegate
source checks to the same module version and preserve their own file exclusions.
golayout needs no sibling module. See the optional
[workspace instructions](development.md#local-workspace) for ecosystem development.
