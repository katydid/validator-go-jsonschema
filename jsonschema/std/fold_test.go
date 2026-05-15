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

import "testing"

func TestMustFold(t *testing.T) {
	xs := []int{1, 2, 3}
	sum := MustFoldA(xs, func(a, b int) int { return a + b })
	if sum != 6 {
		t.Fatalf("expected sum = 6, but got %d", sum)
	}
}
