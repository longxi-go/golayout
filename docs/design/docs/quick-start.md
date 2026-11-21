# Quick start

Requires Go 1.27.0 or later. A workspace is not required.

## Install the CLI

```sh
go install github.com/longxi-go/golayout/cmd/golayout@v0.1.1
```

Add `GOBIN`, or `$(go env GOPATH)/bin` when `GOBIN` is unset, to `PATH`.
To manage the CLI as a module tool instead, run this inside that module:

```sh
go get -tool github.com/longxi-go/golayout/cmd/golayout@v0.1.1
go tool golayout -lang go1.27 path/to/file.go
```

## Check or rewrite files

The installed CLI accepts one or more explicit files:

```sh
golayout -lang go1.27 path/to/file.go
golayout -w -lang go1.27 path/to/file.go
```

The first command is read-only. Differences and invalid input return a nonzero
exit status. `-w` writes changes. From a source checkout, use
`go run ./cmd/golayout` in place of `golayout`.

## Use the library

In an existing Go module, add the formatting package:

```sh
go get github.com/longxi-go/golayout/pkg/format@v0.1.1
```

Save this program as `main.go` and run `go run .`:

```go
package main

import (
	"fmt"
	"log"

	"github.com/longxi-go/golayout/pkg/format"
)

func main() {
	source := []byte("package sample\nfunc Ready() bool { return true }\n")
	formatted, err := format.Source(source, format.Options{
		LangVersion: "go1.27",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(formatted))
}
```

The result contains a multiline `Ready` function. `Source` leaves its input
unchanged. Generators should derive `LangVersion` from the target module's
`go.mod`, handle errors, then validate and persist the result.

## Maintain source notices

Create a plain-text `codequality/HEADER` with the repository's copyright owner
and license notice, then run:

```sh
golayout -header codequality/HEADER path/to/file.go
golayout -w -header codequality/HEADER path/to/file.go
```

The default is read-only. Header mode preserves the code body and rejects
replacing a different owner. It cannot be combined with `-lang`; run formatting
separately. Generators can use `pkg/header.Source` before formatting.

Repository-wide selection and local quality checks are covered by the
[development guide](development.md).
