// Copyright © 2026-present the longxi-go/golayout authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/longxi-go/golayout/pkg/format"
)

func TestRunMatchesLibrary(t *testing.T) {
	input := []byte("package sample\nvar (\nName = \"value\"\n)\nfunc Read() string { return Name }\n")
	path := filepath.Join(t.TempDir(), "sample.go")
	if err := os.WriteFile(path, input, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"-w", path}, io.Discard); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want, err := format.Source(input, format.Options{})
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("CLI differs from library: %v\n%s\n%s", err, got, want)
	}
	if err := Run([]string{"-lang", "go1.27", path}, io.Discard); err != nil {
		t.Fatal(err)
	}
}
