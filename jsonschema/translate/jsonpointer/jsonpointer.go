// Package jsonpointer implements IETF rfc6901
// JSON Pointers are a string syntax for
// identifying a specific value within a JavaScript Object Notation
// (JSON) document [RFC4627].  JSON Pointer is intended to be easily
// expressed in JSON string values as well as Uniform Resource
// Identifier (URI) [RFC3986] fragment identifiers.
//
// this package is intended to work like net/url from the go
// standard library
package jsonpointer

import (
	"fmt"
	"strings"
)

// The ABNF syntax of a JSON Pointer is:
// json-pointer    = *( "/" reference-token )
// reference-token = *( unescaped / escaped )
// unescaped       = %x00-2E / %x30-7D / %x7F-10FFFF
//
//	; %x2F ('/') and %x7E ('~') are excluded from 'unescaped'
//
// escaped         = "~" ( "0" / "1" )
//
//	; representing '~' and '/', respectively
func ParseFragment(str string) ([]string, error) {
	if len(str) == 0 {
		return []string{}, nil
	}

	if str[0] != '/' {
		return nil, fmt.Errorf("non-empty references must begin with a '/' character")
	}
	str = str[1:]

	toks := strings.Split(str, separator)
	for i, t := range toks {
		toks[i] = unescapeToken(t)
	}
	return toks, nil
}

const (
	separator        = "/"
	escapedSeparator = "~1"
	tilde            = "~"
	escapedTilde     = "~0"
)

func unescapeToken(tok string) string {
	tok = strings.Replace(tok, escapedSeparator, separator, -1)
	return strings.Replace(tok, escapedTilde, tilde, -1)
}

func escapeToken(tok string) string {
	tok = strings.Replace(tok, tilde, escapedTilde, -1)
	return strings.Replace(tok, separator, escapedSeparator, -1)
}
