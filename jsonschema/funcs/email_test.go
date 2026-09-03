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

	"katydid.org.za/go/validator-go-jsonschema/jsonschema/funcs/email/lexer"
)

func TestEmail(t *testing.T) {
	var valid = map[string]string{
		"joe.bloggs@example.com":     "a valid e-mail address",
		"te~st@example.com":          "tilde in local part is valid",
		"~test@example.com":          "tilde before local part is valid",
		"test~@example.com":          "tilde after local part is valid",
		"te.s.t@example.com":         "two separated dots inside local part are valid",
		"\"\"@iana.org":              "an empty quoted string in the local part is valid",
		"a@[ipv6:::1]":               "a lowercase IPv6 tag in an address literal is valid",
		"\"joe@bloggs\"@example.com": "a quoted string with a @ in the local part is valid",
		"joe.bloggs@[127.0.0.1]":     "an IPv4-address-literal after the @ is valid",
		"joe.bloggs@[IPv6:::1]":      "an IPv6-address-literal after the @ is valid",
		"\"\\a\"@iana.org":           "a quoted pair in the local part is valid",
		"\"\\\"\"@iana.org":          "an escaped double quote in the local part is valid",
		"\"\\\\\"@iana.org":          "an escaped backslash in the local part is valid",
	}

	var invalid = map[string]string{
		"2962":                                 "an invalid e-mail address",
		".test@example.com":                    "dot before local part is not valid",
		"test.@example.com":                    "dot after local part is not valid",
		"te..st@example.com":                   "two subsequent dots inside local part are not valid",
		"user1@oceania.org, user2@oceania.org": "two email addresses is not valid",
		"\"test\\\u00a9\"@iana.org":            "a non-ASCII character in a quoted pair is not valid",
		"test@iana.org\n":                      "a trailing line feed is not valid",
		"test@iana.org\r":                      "a trailing carriage return is not valid",
		"a@iana.org ":                          "a trailing space is not valid",
		" a@iana.org":                          "a leading space is not valid",
		"test@-iana.org":                       "a domain label starting with a hyphen is not valid",
		"test@iana-.com":                       "a domain label ending with a hyphen is not valid",
		"joe.bloggs@invalid=domain.com":        "an invalid domain",
		"joe bloggs@example.com":               "unquoted space in local part is invalid",
		"a\tb@iana.org":                        "a tab in the local part is not valid",
	}
	lex := lexer.NewLexer([]byte{})
	for email, desc := range valid {
		if !isEmail(lex, []byte(email)) {
			t.Errorf("got false, but expected true for %s: %s", email, desc)
		}
	}
	for email, desc := range invalid {
		if isEmail(lex, []byte(email)) {
			t.Errorf("got true, but expected false for %s: %s", email, desc)
		}
	}
}
