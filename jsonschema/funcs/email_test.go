// Copyright 2026 Walter Schulze
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

package funcs

import (
	"testing"

	"github.com/katydid/parser-go/pool"
	"github.com/katydid/validator-go-jsonschema/jsonschema/funcs/email/lexer"
)

func TestEmail(t *testing.T) {
	pool := pool.New()
	var valid = map[string]string{
		"joe.bloggs@example.com": "a valid e-mail address",
		"te~st@example.com":      "tilde in local part is valid",
		"~test@example.com":      "tilde before local part is valid",
		"test~@example.com":      "tilde after local part is valid",
		"te.s.t@example.com":     "two separated dots inside local part are valid",
	}

	var invalid = map[string]string{
		"2962":               "an invalid e-mail address",
		".test@example.com":  "dot before local part is not valid",
		"test.@example.com":  "dot after local part is not valid",
		"te..st@example.com": "two subsequent dots inside local part are not valid",
	}
	lex := lexer.NewLexer([]byte{})
	for email, desc := range valid {
		if !isEmail(pool, lex, email) {
			t.Fatalf("got false, but expected true for %s: %s", email, desc)
		}
	}
	for email, desc := range invalid {
		if isEmail(pool, lex, email) {
			t.Fatalf("got true, but expected false for %s: %s", email, desc)
		}
	}
}
