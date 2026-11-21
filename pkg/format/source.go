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

package format

import (
	"fmt"
	"go/version"

	"github.com/longxi-go/golayout/internal/layout"
	"mvdan.cc/gofumpt/format"
)

// Source applies the shared policy in memory. Invalid input returns no output;
// callers retain ownership of file discovery, cancellation and persistence.
func Source(input []byte, options Options) ([]byte, error) {
	if options.LangVersion != "" && version.Lang(options.LangVersion) == "" {
		return nil, fmt.Errorf("invalid Go language version %q", options.LangVersion)
	}
	expanded, err := layout.Source(input)
	if err != nil {
		return nil, err
	}
	return format.Source(expanded, format.Options{
		LangVersion: options.LangVersion,
	})
}
