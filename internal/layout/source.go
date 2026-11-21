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
	"go/format"
	"go/parser"
	"go/token"
)

// Source expands structural layout without resolving packages or changing APIs.
// Invalid syntax returns an error; input storage is never modified.
func Source(input []byte) ([]byte, error) {
	positions := token.NewFileSet()
	file, err := parser.ParseFile(positions, "", input, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	editor := editor{
		source:    input,
		positions: positions,
		changes:   make(map[int]change),
	}
	editor.declarations(file)
	ast.Inspect(file, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.FuncDecl:
			editor.function(node.Body)
		case *ast.FuncLit:
			editor.function(node.Body)
		case *ast.StructType:
			editor.fields(node.Fields)
		case *ast.InterfaceType:
			editor.fields(node.Methods)
		case *ast.CompositeLit:
			editor.literal(node)
		}
		return true
	})
	return format.Source(editor.apply())
}
