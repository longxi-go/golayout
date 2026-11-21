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
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"

	"github.com/longxi-go/golayout/pkg/format"
)

func TestSourceExpandsControlBlocksAndClauses(t *testing.T) {
	input, err := os.ReadFile("testdata/blocks.go.txt")
	if err != nil {
		t.Fatal(err)
	}
	got, err := format.Source(input, format.Options{})
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", got, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	ast.Inspect(file, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.BlockStmt:
			if len(node.List) > 0 && fset.Position(node.Lbrace).Line == fset.Position(node.Rbrace).Line {
				t.Errorf("compact executable block at %v", fset.Position(node.Lbrace))
			}
		case *ast.CaseClause:
			assertStatementsBelow(t, fset, node.Colon, node.Body)
		case *ast.CommClause:
			assertStatementsBelow(t, fset, node.Colon, node.Body)
		}
		return true
	})
	if !strings.Contains(string(got), "select {}") {
		t.Fatalf("canonical empty select expanded:\n%s", got)
	}
	again, err := format.Source(got, format.Options{})
	if err != nil || !bytes.Equal(got, again) {
		t.Fatalf("not idempotent: %v\n%s", err, again)
	}
}

func assertStatementsBelow(t *testing.T, fset *token.FileSet, colon token.Pos, body []ast.Stmt) {
	t.Helper()
	previous := fset.Position(colon).Line
	for _, statement := range body {
		line := fset.Position(statement.Pos()).Line
		if line <= previous {
			t.Errorf("compact clause statement at %v", fset.Position(statement.Pos()))
		}
		previous = fset.Position(statement.End()).Line
	}
}
