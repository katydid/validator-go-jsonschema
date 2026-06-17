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

package regexformat

import "strings"

func tryFastPath(expr string) func(s string) bool {
	if fast := tryFastPathAny(expr); fast != nil {
		return fast
	}
	if fast := tryFastPathCharSet(expr); fast != nil {
		return fast
	}
	if fast := tryFastPathPrefix(expr); fast != nil {
		return fast
	}
	return nil
}

// TODO: fast paths to consider building in future
// "^[a-zA-Z0-9\\/_]{1,30}$"
// "^[A-F0-9]{1,32}$"
// "^[a-zA-Z0-9_\\.\\-/|@#]*$"
// "^[a-z][a-z0-9_]+$"
// "^[a-z][a-z0-9-_]{1,63}$"
// "^[1-9][0-9]*$"
// "^:[0-9]+$"
// "^[0-9]+(ns|ms|us|µs|s|m|h)$"
// "^#[0-9a-fA-F]{6}$",
// "^[a-z]{1,2}$"

func tryFastPathAny(expr string) func(s string) bool {
	switch expr {
	case "^.*$", ".*":
		return func(string) bool {
			return true
		}
	case "^.+$", ".+":
		return func(s string) bool {
			return len(s) != 0
		}
	}
	return nil
}

// "^/.*"
func tryFastPathPrefix(expr string) func(s string) bool {
	if expr[0] != '^' {
		return nil
	}
	expr = expr[1:]
	if expr[len(expr)-1] != '*' {
		return nil
	}
	expr = expr[:len(expr)-1]
	if expr[len(expr)-1] != '.' {
		return nil
	}
	expr = expr[:len(expr)-1]
	for i := 0; i < len(expr); i++ {
		if !isUnreservedAscii(expr[i]) {
			return nil
		}
	}
	prefix := expr
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}

func isUnreservedAscii(c byte) bool {
	if c >= 'A' && c <= 'Z' {
		return true
	} else if c >= 'a' && c <= 'z' {
		return true
	} else if c >= '0' && c <= '9' {
		return true
	} else if c == '_' || c == '/' || c == '-' {
		return true
	}
	return false
}

func tryFastPathCharSet(expr string) func(s string) bool {
	if expr[0] != '^' {
		return nil
	}
	expr = expr[1:]
	if expr[0] != '[' {
		return nil
	}
	expr = expr[1:]
	if expr[len(expr)-1] != '$' {
		return nil
	}
	expr = expr[:len(expr)-1]
	plus := false
	star := false
	if expr[len(expr)-1] == '+' {
		plus = true
		expr = expr[:len(expr)-1]
	} else if expr[len(expr)-1] == '*' {
		plus = true
		expr = expr[:len(expr)-1]
	}
	offset, set := getCharSet(expr)
	if offset == -1 {
		return nil
	}
	if offset != len(expr)-1 {
		return nil
	}
	if expr[offset] != ']' {
		return nil
	}
	if plus {
		return matchAlpha1(set)
	}
	if star {
		return matchAlpha(set)
	}
	return nil
}

func getCharSet(expr string) (int, [256]byte) {
	set := [256]byte{}
	var prev byte
	i := 0
	for i < len(expr) {
		c := expr[i]
		if isUnreservedInsideRange(c) {
			set[c] = 1
		} else if c == '-' {
			next := expr[i+1]
			if !isUnreservedInsideRange(next) {
				return -1, set
			}
			for j := prev; j <= next; j++ {
				set[j] = 1
			}
		} else if c == ']' {
			return i, set
		} else {
			return -1, set
		}
		prev = c
		i++
	}
	return -1, set
}

func isUnreservedInsideRange(c byte) bool {
	if c >= 'A' && c <= 'Z' {
		return true
	} else if c >= 'a' && c <= 'z' {
		return true
	} else if c >= '0' && c <= '9' {
		return true
	} else if c == '_' || c == '/' || c == '.' {
		return true
	}
	return false
}

func matchAlpha1(alphabet [256]byte) func(s string) bool {
	return func(s string) bool {
		if len(s) == 0 {
			return false
		}
		for _, c := range s {
			if alphabet[c] != 1 {
				return false
			}
		}
		return true
	}
}

func matchAlpha(alphabet [256]byte) func(s string) bool {
	return func(s string) bool {
		for _, c := range s {
			if alphabet[c] != 1 {
				return false
			}
		}
		return true
	}
}
