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
	"go/token"
	"testing"

	golayout "github.com/longxi-go/golayout/internal/layout"
)

func FuzzSource(f *testing.F) {
	for _, source := range []string{
		"package sample\ntype ()\n",
		"package sample\ntype Config struct{ Value int };func Read() int{return 1}\n",
		"package sample\nvar Values = []Value{{Name: \"n\" /* trailing */}}\n",
		"package sample\nfunc Read(){ f := func() { }; f() }\n",
	} {
		f.Add(source)
	}
	f.Fuzz(func(t *testing.T, source string) {
		if _, err := parser.ParseFile(token.NewFileSet(), "", source, parser.ParseComments); err != nil {
			return
		}
		formatted, err := golayout.Source([]byte(source))
		if err != nil {
			t.Fatalf("valid Go became invalid: %v", err)
		}
		again, err := golayout.Source(formatted)
		if err != nil || !bytes.Equal(formatted, again) {
			t.Fatalf("non-idempotent output: %v\n%s\n%s", err, formatted, again)
		}
	})
}
