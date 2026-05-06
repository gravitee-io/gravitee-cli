// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmdutil

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// ResolvePassword returns the password to use, applying this precedence:
//  1. --password-stdin: read a single line from stdin (trailing newline stripped).
//  2. --password flag (cleartext): used as-is. Caller should print a deprecation hint.
//  3. interactive: prompt on stderr if stdin is a terminal.
//
// Returns an error if no source is available and stdin is not a TTY.
//
// The cleartext --password path is intentionally last-resort: it's visible in
// `ps`, shell history, and process audit logs.
func ResolvePassword(flagValue string, fromStdin bool, prompt string, stdin io.Reader, stderr io.Writer) (string, error) {
	if fromStdin && flagValue != "" {
		return "", errors.New("--password and --password-stdin are mutually exclusive")
	}

	if fromStdin {
		return readPasswordFromStdin(stdin)
	}

	if flagValue != "" {
		return flagValue, nil
	}

	return promptPassword(prompt, stderr)
}

func readPasswordFromStdin(stdin io.Reader) (string, error) {
	scanner := bufio.NewScanner(stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("read password from stdin: %w", err)
		}
		return "", errors.New("no password provided on stdin")
	}
	pw := strings.TrimRight(scanner.Text(), "\r\n")
	if pw == "" {
		return "", errors.New("empty password on stdin")
	}
	return pw, nil
}

func promptPassword(prompt string, stderr io.Writer) (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return "", errors.New("no password provided: pass --password-stdin or run interactively")
	}
	if prompt == "" {
		prompt = "Password: "
	}
	fmt.Fprint(stderr, prompt)
	pw, err := term.ReadPassword(fd)
	fmt.Fprintln(stderr)
	if err != nil {
		return "", fmt.Errorf("read password: %w", err)
	}
	if len(pw) == 0 {
		return "", errors.New("empty password")
	}
	return string(pw), nil
}
