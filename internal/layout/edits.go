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
	"bytes"
	"go/token"
	"slices"
)

// Edits address original bytes, so nested nodes never invalidate one another's
// positions. A comma precedes newlines when both insertions share an address.
type change struct {
	comma bool
	lines int
}

type editor struct {
	source    []byte
	positions *token.FileSet
	changes   map[int]change
}

func (e *editor) insert(offset, lines int, comma bool) {
	change := e.changes[offset]
	change.lines = max(change.lines, lines)
	change.comma = change.comma || comma
	e.changes[offset] = change
}

func (e *editor) offset(pos token.Pos) int {
	return e.positions.PositionFor(pos, false).Offset
}

func (e *editor) separate(pos token.Pos, blank bool) {
	offset := e.offset(pos)
	start := bytes.LastIndexByte(e.source[:offset], '\n') + 1
	if len(bytes.TrimSpace(e.source[start:offset])) > 0 {
		lines := 1
		if blank {
			lines = 2
		}
		e.insert(offset, lines, false)
	} else if blank && start > 0 {
		previous := bytes.LastIndexByte(e.source[:start-1], '\n') + 1
		if len(bytes.TrimSpace(e.source[previous:start-1])) > 0 {
			e.insert(start, 1, false)
		}
	}
}

func (e *editor) apply() []byte {
	offsets := make([]int, 0, len(e.changes))
	for offset := range e.changes {
		offsets = append(offsets, offset)
	}
	slices.Sort(offsets)
	var output bytes.Buffer
	previous := 0
	for _, offset := range offsets {
		output.Write(e.source[previous:offset])
		change := e.changes[offset]
		if change.comma {
			output.WriteByte(',')
		}
		for range change.lines {
			output.WriteByte('\n')
		}
		previous = offset
	}
	output.Write(e.source[previous:])
	return output.Bytes()
}
