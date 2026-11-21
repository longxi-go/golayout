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

package header

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type removal struct {
	start int
	end   int
}

// Source inserts or refreshes a source notice without formatting the code body.
// Existing notices must identify the template's owner; a different owner is an
// error, so repository-wide refreshes cannot silently relabel third-party code.
// Generated markers stay first, preserving generator ownership checks. Invalid
// Go or an ambiguous notice returns no output, leaving persistence to callers.
func Source(source, notice []byte) ([]byte, error) {
	rendered, owner, err := prepareNotice(notice)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "source.go", source, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	var comments []*ast.Comment
	for _, group := range file.Comments {
		if group.Pos() < file.Package {
			comments = append(comments, group.List...)
		}
	}
	var edits []removal
	var generated string
	for i := 0; i < len(comments); i++ {
		comment := comments[i]
		text := commentText(comment)
		start := fset.Position(comment.Pos()).Offset
		end := commentEnd(source, start)
		if strings.HasPrefix(comment.Text, "// Code generated ") && strings.HasSuffix(comment.Text, " DO NOT EDIT.") {
			if generated != "" {
				return nil, fmt.Errorf("multiple generated markers in source preamble")
			}
			generated = comment.Text + "\n\n"
			edits = append(edits, removal{
				start: start,
				end:   end,
			})
			continue
		}
		if !strings.HasPrefix(text, "Copyright ") {
			continue
		}
		last := i
		if strings.HasPrefix(comment.Text, "//") {
			last, err = lineNoticeEnd(comments, i)
			if err != nil {
				return nil, err
			}
		}
		for _, part := range comments[i : last+1] {
			if err := checkOwners(commentText(part), owner); err != nil {
				return nil, err
			}
		}
		end = commentEnd(source, fset.Position(comments[last].Pos()).Offset)
		i = last
		edits = append(edits, removal{
			start: start,
			end:   end,
		})
	}
	remaining := bytes.Clone(source)
	for i := len(edits) - 1; i >= 0; i-- {
		edit := edits[i]
		remaining = append(remaining[:edit.start], remaining[edit.end:]...)
	}
	// A UTF-8 BOM is allowed only at the beginning of a Go source file.
	var result []byte
	if bytes.HasPrefix(remaining, []byte{0xef, 0xbb, 0xbf}) {
		result = append(result, remaining[:3]...)
		remaining = remaining[3:]
	}
	result = append(result, generated...)
	result = append(result, rendered...)
	return append(result, bytes.TrimLeft(remaining, " \t\r\n")...), nil
}
