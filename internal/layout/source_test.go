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

package layout_test

import (
	"bytes"
	"go/parser"
	"go/scanner"
	"go/token"
	"slices"
	"strings"
	"testing"

	golayout "github.com/longxi-go/golayout/internal/layout"
)

func TestStructuralLayout(t *testing.T) {
	for _, test := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty declaration group",
			input: "type ()\n",
			want:  "type ()\n",
		},
		{
			name:  "named and anonymous structs",
			input: "type Config struct{ Port int; Name string }\nvar Inline struct{ Value int }\n",
			want:  "type Config struct {\n\tPort int\n\tName string\n}\n\nvar Inline struct {\n\tValue int\n}\n",
		},
		{
			name:  "methods and empty hooks",
			input: "type Config struct{}\nfunc (Config) Port() int { return 8080 }\nfunc (Config) Ready() {}\n",
			want:  "type Config struct{}\n\nfunc (Config) Port() int {\n\treturn 8080\n}\n\nfunc (Config) Ready() {}\n",
		},
		{
			name:  "closure and semicolons",
			input: "var Read = func() int { n := 0; for i := 0; i < 2; i++ { n += i }; return n }\n",
			want:  "var Read = func() int {\n\tn := 0\n\tfor i := 0; i < 2; i++ {\n\t\tn += i\n\t}\n\treturn n\n}\n",
		},
		{
			name:  "single field and generic literals",
			input: "var Config = Box[int]{Value: 1}\nvar Empty = Box[int]{}\n",
			want:  "var Config = Box[int]{\n\tValue: 1,\n}\n\nvar Empty = Box[int]{}\n",
		},
		{
			name:  "nested and inferred literals",
			input: "var Configs = []Config{{Port: 1, Name: \"one\"}}\n",
			want:  "var Configs = []Config{\n\t{\n\t\tPort: 1,\n\t\tName: \"one\",\n\t},\n}\n",
		},
		{
			name:  "grouped interfaces and types",
			input: "type (\nA interface{ Read() int }\nB struct{ Name string }\n)\n",
			want:  "type (\n\tA interface {\n\t\tRead() int\n\t}\n\n\tB struct {\n\t\tName string\n\t}\n)\n",
		},
		{
			name:  "same line declarations",
			input: "type A struct{}; type B struct{}; func Read() int { return 1 }\n",
			want:  "type A struct{}\n\ntype B struct{}\n\nfunc Read() int {\n\treturn 1\n}\n",
		},
		{
			name:  "existing multiline literal fields",
			input: "var Config = Box{\nA: 1, B: 2,\n}\n",
			want:  "var Config = Box{\n\tA: 1,\n\tB: 2,\n}\n",
		},
		{
			name:  "empty bodies and simple collections",
			input: "type Empty struct{}\nvar Any interface{}\nvar Names = []string{\"a\", \"b\"}\n",
			want:  "type Empty struct{}\n\nvar Any interface{}\n\nvar Names = []string{\"a\", \"b\"}\n",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			input := []byte("package fixture\n\n" + test.input)
			original := bytes.Clone(input)
			got, err := golayout.Source(input)
			if err != nil {
				t.Fatal(err)
			}
			if want := "package fixture\n\n" + test.want; string(got) != want {
				t.Fatalf("want:\n%s\ngot:\n%s", want, got)
			}
			if !bytes.Equal(input, original) || !slices.Equal(tokens(t, original), tokens(t, got)) {
				t.Fatal("formatting changed source tokens or input storage")
			}
			again, err := golayout.Source(got)
			if err != nil || !bytes.Equal(got, again) {
				t.Fatalf("not idempotent: %v\n%s", err, again)
			}
		})
	}
}

func TestCommentsTagsAndRawStrings(t *testing.T) {
	input := "//go:build tooling\n\npackage fixture\n\n// Config retains the marker.\n// +ioc:component\ntype Config struct{ Value string `property:\"a;b\"` }\nfunc Read() string { /* retained */ return `line;one\nline two` }\n"
	got, err := golayout.Source([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{"//go:build tooling\n\npackage", "// +ioc:component\ntype Config", "`property:\"a;b\"`", "/* retained */", "`line;one\nline two`"} {
		if !strings.Contains(string(got), fragment) {
			t.Errorf("lost %q", fragment)
		}
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "result.go", got, parser.ParseComments); err != nil {
		t.Fatal(err)
	}
}

func TestInvalidSyntax(t *testing.T) {
	if _, err := golayout.Source([]byte("package fixture\nfunc Broken(")); err == nil {
		t.Fatal("invalid syntax accepted")
	}
}

func tokens(t *testing.T, source []byte) []string {
	t.Helper()
	var scan scanner.Scanner
	file := token.NewFileSet().AddFile("source.go", -1, len(source))
	scan.Init(file, source, nil, 0)
	var result []string
	for {
		_, tok, literal := scan.Scan()
		if tok == token.EOF {
			return result
		}
		if tok != token.SEMICOLON && tok != token.COMMA {
			result = append(result, tok.String()+":"+literal)
		}
	}
}
