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
	"strings"
	"testing"
)

func TestCheckWriteAndIdempotence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "with space.go")
	before := "package sample\n\ntype Config struct{ Value int }\n"
	if err := os.WriteFile(path, []byte(before), 0o640); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := Run([]string{path}, &output); err == nil || !strings.Contains(output.String(), path) {
		t.Fatalf("check accepted compact source: %v %q", err, output.String())
	}
	if data, err := os.ReadFile(path); err != nil || string(data) != before {
		t.Fatal("check mutated source", err)
	}
	if err := Run([]string{"-w", path}, io.Discard); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(data), "struct {\n\tValue int\n}") {
		t.Fatalf("layout not written: %v %s", err, data)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0o640 {
		t.Fatal("permissions changed", err)
	}
	if err := Run([]string{path}, io.Discard); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"-w", path}, io.Discard); err != nil {
		t.Fatal(err)
	}
	again, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(data, again) {
		t.Fatal("write not idempotent", err)
	}
}

func TestInvalidInputDoesNotWriteEarlierFiles(t *testing.T) {
	dir := t.TempDir()
	valid, broken := filepath.Join(dir, "valid.go"), filepath.Join(dir, "broken.go")
	before := "package sample\nfunc Read() int{return 1}\n"
	if err := os.WriteFile(valid, []byte(before), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(broken, []byte("package sample\nfunc Broken("), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"-w", valid, broken}, io.Discard); err == nil {
		t.Fatal("invalid syntax accepted")
	}
	if data, err := os.ReadFile(valid); err != nil || string(data) != before {
		t.Fatal("earlier source changed after invalid input", err)
	}
	for _, args := range [][]string{
		nil,
		{
			"-unknown",
		},
		{
			filepath.Join(dir, "absent.go"),
		},
		{
			dir,
		},
	} {
		if err := Run(args, io.Discard); err == nil {
			t.Fatalf("invalid invocation accepted: %v", args)
		}
	}
}
