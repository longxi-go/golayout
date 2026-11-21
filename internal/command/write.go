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
	"fmt"
	"os"
	"path/filepath"
)

// Replace each file only after complete writes and close; partial I/O must never
// truncate source. This does not promise a transaction across multiple files.
func (f formattedFile) write() error {
	file, err := os.CreateTemp(filepath.Dir(f.path), ".golayout-*")
	if err != nil {
		return fmt.Errorf("%s: %w", f.path, err)
	}
	name := file.Name()
	defer func() {
		_ = os.Remove(name)
	}()
	if _, err = file.Write(f.data); err != nil {
		_ = file.Close()
		return fmt.Errorf("%s: %w", f.path, err)
	}
	if err = file.Close(); err != nil {
		return fmt.Errorf("%s: %w", f.path, err)
	}
	if err = os.Chmod(name, f.mode); err != nil {
		return fmt.Errorf("%s: %w", f.path, err)
	}
	if err = os.Rename(name, f.path); err != nil {
		return fmt.Errorf("%s: %w", f.path, err)
	}
	return nil
}
