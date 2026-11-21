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
)

func TestHeaderCheckWriteAndPreflight(t *testing.T) {
	dir := t.TempDir()
	template := filepath.Join(dir, "HEADER")
	path := filepath.Join(dir, "source with space.go")
	if err := os.WriteFile(template, []byte("Copyright © 2026 the example authors.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := []byte("package example\n\nfunc Value() int { return 1 }\n")
	if err := os.WriteFile(path, before, 0o640); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := Run([]string{"-header", template, path}, &output); err == nil || !bytes.Contains(output.Bytes(), []byte(path)) {
		t.Fatalf("header drift not reported: %v %q", err, output.String())
	}
	if got, err := os.ReadFile(path); err != nil || !bytes.Equal(got, before) {
		t.Fatal("check changed source", err)
	}
	missing := filepath.Join(dir, "missing.go")
	if err := Run([]string{"-w", "-header", template, path, missing}, io.Discard); err == nil {
		t.Fatal("invalid file accepted")
	}
	if got, err := os.ReadFile(path); err != nil || !bytes.Equal(got, before) {
		t.Fatal("preflight failure changed earlier file", err)
	}
	if err := Run([]string{"-w", "-header", template, path}, io.Discard); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil || !bytes.HasPrefix(got, []byte("// Copyright")) || !bytes.HasSuffix(got, before) {
		t.Fatalf("header mode changed body formatting: %v\n%s", err, got)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0o640 {
		t.Fatal("file mode changed", err)
	}
	if err := Run([]string{"-header", template, path}, io.Discard); err != nil {
		t.Fatal(err)
	}
}
