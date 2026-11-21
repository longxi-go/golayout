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
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/longxi-go/golayout/pkg/format"
	"github.com/longxi-go/golayout/pkg/header"
)

// Run checks explicit files by default and writes only when requested with -w.
// Differences and operational failures return errors for the CLI exit status.
func Run(args []string, output io.Writer) error {
	flags := flag.NewFlagSet("golayout", flag.ContinueOnError)
	flags.SetOutput(output)
	write := flags.Bool("w", false, "rewrite files instead of checking them")
	langVersion := flags.String("lang", "", "Go language version passed to gofumpt")
	headerPath := flags.String("header", "", "check or refresh source notices from a plain-text template, without formatting")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if flags.NArg() == 0 {
		return errors.New("usage: golayout [-w] [-lang version | -header template] <file.go> [file.go]")
	}
	files, err := prepareRequested(flags.Args(), *langVersion, *headerPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if *write {
			if err := file.write(); err != nil {
				return err
			}
		} else if _, err := fmt.Fprintln(output, file.path); err != nil {
			return err
		}
	}
	if !*write && len(files) > 0 {
		if *headerPath != "" {
			return errors.New("source headers differ; run make license-headers")
		}
		return errors.New("go layout differs; run make format")
	}
	return nil
}

func prepareRequested(paths []string, language, template string) ([]formattedFile, error) {
	if template == "" {
		return prepare(paths, format.Options{
			LangVersion: language,
		})
	}
	if language != "" {
		return nil, errors.New("-header and -lang are separate operations")
	}
	notice, err := os.ReadFile(template)
	if err != nil {
		return nil, err
	}
	return prepareWith(paths, func(source []byte) ([]byte, error) {
		return header.Source(source, notice)
	})
}
