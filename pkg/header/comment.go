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
	"strings"
)

// AST comment text omits carriage returns, so locate raw delimiters rather than
// deriving byte offsets from its length when handling existing CRLF block headers.
func commentEnd(source []byte, start int) int {
	if bytes.HasPrefix(source[start:], []byte("/*")) {
		return start + bytes.Index(source[start:], []byte("*/")) + 2
	}
	if newline := bytes.IndexByte(source[start:], '\n'); newline >= 0 {
		return start + newline
	}
	return len(source)
}

func commentText(comment *ast.Comment) string {
	text := strings.TrimPrefix(comment.Text, "//")
	if strings.HasPrefix(text, "/*") {
		text = strings.TrimSuffix(strings.TrimPrefix(text, "/*"), "*/")
		lines := strings.Split(text, "\n")
		for i, line := range lines {
			lines[i] = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "*"))
		}
		text = strings.Join(lines, "\n")
	}
	return strings.TrimSpace(text)
}

func lineNoticeEnd(comments []*ast.Comment, start int) (int, error) {
	next := start + 1
	for next < len(comments) && commentText(comments[next]) == "" {
		next++
	}
	if next == len(comments) {
		return start, nil
	}
	text := commentText(comments[next])
	if strings.HasPrefix(text, "SPDX-License-Identifier:") {
		return next, nil
	}
	if !strings.HasPrefix(text, "Licensed under the Apache License,") {
		return start, nil
	}
	for i := next; i < len(comments); i++ {
		if commentText(comments[i]) == "limitations under the License." {
			return i, nil
		}
	}
	return 0, fmt.Errorf("incomplete Apache source header; refusing to guess its boundary")
}
