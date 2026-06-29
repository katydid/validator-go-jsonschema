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

import (
	"strconv"
	"strings"
)

// TODO: fast paths to consider building in future
// TODO: consider including `\\/` and `\\.` and `\\-` and `/` and `|` and `@` and `#` in charset
// TODO: length min and max "^[charset]{1,30}$"
// TODO: support a prefix before a charset "^prefix[charset]+$"
// TODO: support times "^[0-9]+(ns|ms|us|µs|s|m|h)$"

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
	if fast := tryFastPathLength(expr); fast != nil {
		return fast
	}
	return nil
}

func tryFastPathLength(expr string) func(s string) bool {
	if expr[0] != '^' {
		return nil
	}
	expr = expr[1:]
	if expr[len(expr)-1] != '$' {
		return nil
	}
	expr = expr[:len(expr)-1]
	if expr[0] != '.' {
		return nil
	}
	expr = expr[1:]
	if expr[0] != '{' {
		return nil
	}
	expr = expr[1:]
	if expr[len(expr)-1] != '}' {
		return nil
	}
	expr = expr[:len(expr)-1]
	ss := strings.Split(expr, ",")
	if len(ss) == 1 {
		exact, err := strconv.Atoi(ss[0])
		if err != nil {
			return nil
		}
		return func(s string) bool {
			l := 0
			if len(s) < exact {
				return false
			}
			for range s {
				l++
				if l > exact {
					return false
				}
			}
			return l == exact
		}
	}
	if len(ss) == 2 {
		min, err := strconv.Atoi(ss[0])
		if err != nil {
			return nil
		}
		max, err := strconv.Atoi(ss[1])
		if err != nil {
			return nil
		}
		return func(s string) bool {
			l := 0
			if len(s) < min {
				return false
			}
			for range s {
				l++
				if l > max {
					return false
				}
			}
			return l >= min && l <= max
		}
	}
	return nil
}

func tryFastPathAny(expr string) func(s string) bool {
	switch expr {
	case "^.*$", ".*":
		return func(string) bool {
			return true
		}
	case "^.+$", ".+", ".", "(.+)":
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
	// ignore .*
	if expr[len(expr)-1] == '*' {
		expr = expr[:len(expr)-1]
		if expr[len(expr)-1] == '.' {
			expr = expr[:len(expr)-1]
		} else {
			return nil
		}
	}
	for i := 0; i < len(expr); i++ {
		if !isUnreservedAscii(expr[i]) {
			if i == len(expr)-1 && expr[i] == '-' {
				continue
			}
			// if '-' is last character it is not reserved
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
	justprefix := true
	if expr[len(expr)-1] == '$' {
		justprefix = false
		expr = expr[:len(expr)-1]
	}
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
	if justprefix && !star && !plus {
		return matchFirstCharSet(set)
	}
	if plus {
		return matchCharSet1(set)
	}
	if star {
		return matchCharSet(set)
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
			i++
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
	} else if c == '#' || c == '@' || c == '$' {
		return true
	}
	return false
}

func matchCharSet1(alphabet [256]byte) func(s string) bool {
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

func matchCharSet(alphabet [256]byte) func(s string) bool {
	return func(s string) bool {
		for _, c := range s {
			if alphabet[c] != 1 {
				return false
			}
		}
		return true
	}
}

func matchFirstCharSet(alphabet [256]byte) func(s string) bool {
	return func(s string) bool {
		if len(s) == 0 {
			return false
		}
		return alphabet[s[0]] == 1
	}
}
