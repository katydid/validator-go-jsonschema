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

package std

import (
	"testing"

	"github.com/katydid/validator-go-jsonschema/jsonschema/std/testutil"
)

func TestComplementarySubsets(t *testing.T) {
	want := []struct {
		Left  []int
		Right []int
	}{
		{Left: []int{1}, Right: []int{2}},
		{Left: []int{2}, Right: []int{1}},
		{Left: []int{1, 2}, Right: []int{}},
	}
	got := ComplementarySubsets([]int{1, 2})
	testutil.Expect(t, "", want, got)
}

func TestSubsets(t *testing.T) {
	want := [][]int{
		{1},
		{2},
		{1, 2},
	}
	got := Subsets([]int{1, 2})
	testutil.Expect(t, "", want, got)
}
