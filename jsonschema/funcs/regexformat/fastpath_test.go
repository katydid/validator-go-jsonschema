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

import "testing"

func TestFastPathAny(t *testing.T) {
	expr := "^.*$"
	match := tryFastPath(expr)
	if match == nil {
		t.Fatal("expected fast path")
	}
	if !match("abc") {
		t.Fatal()
	}
}

func TestFastPathPrefix(t *testing.T) {
	expr := "^/.*"
	match := tryFastPath(expr)
	if match == nil {
		t.Fatal("expected fast path")
	}
	tests := map[string]bool{
		"/a":       true,
		"/abc":     true,
		"/abc/def": true,
		"abc":      false,
		"":         false,
	}
	for test, want := range tests {
		t.Run(test, func(t *testing.T) {
			if got := match(test); got != want {
				t.Fatal()
			}
		})
	}
}

func TestFastPathCharSet(t *testing.T) {
	expr := "^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]+$"
	match := tryFastPath(expr)
	if match == nil {
		t.Fatal("expected fast path")
	}
	tests := map[string]bool{
		"123": true,
		"ABC": true,
		"abc": true,
		"9zZ": true,
		"":    false,
		"210": false,
		"012": false,
	}
	for test, want := range tests {
		t.Run(test, func(t *testing.T) {
			if got := match(test); got != want {
				t.Fatal()
			}
		})
	}
}

func TestFastPathRangeCharSet(t *testing.T) {
	expr := "^[0123456789A-Fa-f]+$"
	match := tryFastPath(expr)
	if match == nil {
		t.Fatal("expected fast path")
	}
	tests := map[string]bool{
		"0123456789": true,
		"ABCDEF":     true,
		"fedcba":     true,
		"":           false,
		"ABCDZ":      false,
	}
	for test, want := range tests {
		t.Run(test, func(t *testing.T) {
			if got := match(test); got != want {
				t.Fatal()
			}
		})
	}
}
