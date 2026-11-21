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
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var copyrightPattern = regexp.MustCompile(`(?i)^Copyright\s+(?:(?:©|\(c\))\s+)?[0-9]{4}(?:\s*-\s*(?:[0-9]{4}|present))?\s+(.+)$`)

func copyrightOwner(line string) (string, error) {
	match := copyrightPattern.FindStringSubmatch(strings.TrimSpace(line))
	if len(match) != 2 {
		return "", fmt.Errorf("unsupported copyright line %q", line)
	}
	return strings.TrimSuffix(match[1], " All rights reserved."), nil
}

func checkOwners(text, owner string) error {
	for _, line := range strings.Split(text, "\n") {
		if !strings.HasPrefix(line, "Copyright ") {
			continue
		}
		existingOwner, err := copyrightOwner(line)
		if err != nil {
			return err
		}
		if existingOwner != owner {
			return fmt.Errorf("refusing to replace copyright owned by %q (template owner %q)", existingOwner, owner)
		}
	}
	return nil
}

func prepareNotice(notice []byte) (string, string, error) {
	if !utf8.Valid(notice) || strings.ContainsRune(string(notice), '\x00') {
		return "", "", fmt.Errorf("header template must be UTF-8 text without NUL bytes")
	}
	text := strings.TrimSpace(strings.ReplaceAll(string(notice), "\r\n", "\n"))
	lines := strings.Split(text, "\n")
	owner, err := copyrightOwner(lines[0])
	if err != nil {
		return "", "", fmt.Errorf("header template must start with a plain-text copyright line: %w", err)
	}
	for i, line := range lines {
		line = strings.TrimRight(line, " \t\r")
		if line == "" {
			lines[i] = "//"
		} else {
			lines[i] = "// " + line
		}
	}
	return strings.Join(lines, "\n") + "\n\n", owner, nil
}
