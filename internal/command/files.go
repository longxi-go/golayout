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

package command

import (
	"bytes"
	"fmt"
	"os"

	"github.com/longxi-go/golayout/pkg/format"
)

type formattedFile struct {
	path string
	data []byte
	mode os.FileMode
}

type sourceTransform func([]byte) ([]byte, error)

// Prepare every file before writing any, so a syntax failure cannot leave an
// earlier file rewritten. File ownership remains with the caller's inventory.
func prepare(paths []string, options format.Options) ([]formattedFile, error) {
	return prepareWith(paths, func(source []byte) ([]byte, error) {
		return format.Source(source, options)
	})
}

func prepareWith(paths []string, transform sourceTransform) ([]formattedFile, error) {
	var files []formattedFile
	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("%s: expected a regular source file", path)
		}
		original, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		formatted, err := transform(original)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if !bytes.Equal(original, formatted) {
			files = append(files, formattedFile{
				path: path,
				data: formatted,
				mode: info.Mode().Perm(),
			})
		}
	}
	return files, nil
}
