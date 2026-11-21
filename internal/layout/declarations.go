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
	"go/token"
)

func (e *editor) declarations(file *ast.File) {
	for _, declaration := range file.Decls {
		start := declaration.Pos()
		switch declaration := declaration.(type) {
		case *ast.FuncDecl:
			if declaration.Doc != nil {
				start = declaration.Doc.Pos()
			}
		case *ast.GenDecl:
			if declaration.Doc != nil {
				start = declaration.Doc.Pos()
			}
			if declaration.Tok == token.TYPE && len(declaration.Specs) > 1 {
				for _, spec := range declaration.Specs[1:] {
					typ := spec.(*ast.TypeSpec)
					pos := typ.Pos()
					if typ.Doc != nil {
						pos = typ.Doc.Pos()
					}
					e.separate(pos, true)
				}
			}
		}
		e.separate(start, true)
	}
}
