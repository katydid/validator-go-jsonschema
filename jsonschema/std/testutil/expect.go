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

package testutil

import "testing"

func Expect[A comparable](t *testing.T, desc string, want A, got A) {
	t.Helper()
	if got != want {
		t.Fatalf("%s want %v got %v", desc, want, got)
	}
}

func ExpectErr[A comparable](t *testing.T, desc string, want A, got A) {
	t.Helper()
	if got != want {
		t.Errorf("%s want %v got %v", desc, want, got)
	}
}
