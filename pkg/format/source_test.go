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

package format_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/longxi-go/golayout/pkg/format"
)

func TestSourceComposesLayoutAndGofumpt(t *testing.T) {
	input := []byte("package sample\n\nvar (\nName = \"value\"\n)\n\ntype Config struct{ Port int }\nfunc Read() int { return Config{Port: 1}.Port }\n")
	before := bytes.Clone(input)
	got, err := format.Source(input, format.Options{
		LangVersion: "go1.27",
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{"var Name = \"value\"", "struct {\n\tPort int\n}", "func Read() int {\n", "Config{\n\t\tPort: 1,\n\t}"} {
		if !strings.Contains(string(got), fragment) {
			t.Errorf("missing complete formatting %q:\n%s", fragment, got)
		}
	}
	if !bytes.Equal(input, before) {
		t.Fatal("input storage changed")
	}
	again, err := format.Source(got, format.Options{
		LangVersion: "go1.27",
	})
	if err != nil || !bytes.Equal(got, again) {
		t.Fatalf("formatting is not idempotent: %v", err)
	}
}

func TestSourceRejectsInvalidSyntax(t *testing.T) {
	output, err := format.Source([]byte("package broken\nfunc ("), format.Options{})
	if err == nil || output != nil {
		t.Fatalf("invalid syntax accepted: %s, %v", output, err)
	}
}

func TestSourceRejectsInvalidLanguageVersion(t *testing.T) {
	output, err := format.Source([]byte("package sample\n"), format.Options{
		LangVersion: "1.27",
	})
	if err == nil || output != nil {
		t.Fatalf("invalid language version accepted: %s, %v", output, err)
	}
}
