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

package layout

import (
	"go/ast"
	"go/scanner"
	"go/token"
)

func (e *editor) fields(fields *ast.FieldList) {
	if len(fields.List) == 0 {
		return
	}
	for _, field := range fields.List {
		start := field.Pos()
		if field.Doc != nil {
			start = field.Doc.Pos()
		}
		e.separate(start, false)
	}
	e.separate(fields.Closing, false)
}

func (e *editor) literal(literal *ast.CompositeLit) {
	if len(literal.Elts) == 0 || !expandedLiteral(literal) {
		return
	}
	for _, element := range literal.Elts {
		e.separate(element.Pos(), false)
	}
	e.separate(literal.Rbrace, false)
	last := literal.Elts[len(literal.Elts)-1].End()
	// A newly multiline literal needs its trailing comma before any trailing
	// comment. Existing commas are recognized lexically, without inspecting text.
	var scan scanner.Scanner
	tail := e.source[e.offset(last):e.offset(literal.Rbrace)]
	file := token.NewFileSet().AddFile("", -1, len(tail))
	scan.Init(file, tail, nil, 0)
	for {
		_, tok, _ := scan.Scan()
		if tok == token.COMMA {
			return
		}
		if tok == token.EOF {
			e.insert(e.offset(last), 0, true)
			return
		}
	}
}

func expandedLiteral(literal *ast.CompositeLit) bool {
	switch literal.Type.(type) {
	case *ast.ArrayType, *ast.MapType:
		for _, element := range literal.Elts {
			if entry, ok := element.(*ast.KeyValueExpr); ok {
				element = entry.Value
			}
			switch element.(type) {
			case *ast.CompositeLit, *ast.FuncLit:
				return true
			}
		}
		return false
	default:
		// Named/generic and inferred literals cannot be classified precisely
		// without loading types. Expand them consistently, including one element.
		return true
	}
}
