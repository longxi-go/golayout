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

func (e *editor) function(body *ast.BlockStmt) {
	if body == nil || len(body.List) == 0 || e.positions.PositionFor(body.Lbrace, false).Line != e.positions.PositionFor(body.Rbrace, false).Line {
		return
	}
	left, right := e.offset(body.Lbrace), e.offset(body.Rbrace)
	e.insert(left+1, 1, false)
	e.insert(right, 1, false)
	// Lexical scanning distinguishes statement separators from semicolons inside
	// comments and quoted/raw strings. Retain separators required by for/if heads.
	var scan scanner.Scanner
	file := token.NewFileSet().AddFile("", -1, right-left)
	scan.Init(file, e.source[left:right], nil, scanner.ScanComments)
	for {
		pos, tok, literal := scan.Scan()
		if tok == token.EOF {
			return
		}
		if tok == token.SEMICOLON && literal == ";" {
			e.insert(left+file.Offset(pos)+1, 1, false)
		}
	}
}
