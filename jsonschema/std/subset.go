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

import "github.com/katydid/validator-go/validator/sets"

func ComplementarySubsets[A any](xs []A) []struct {
	Left  []A
	Right []A
} {
	size := len(xs)
	if size == 0 {
		return nil
	}
	max := sets.NewBits(size)
	for i := 0; i < size; i++ {
		max.Set(i, true)
	}
	current := sets.NewBits(len(xs))
	// do not include empty set on the Left side.
	current = current.Inc()

	combos := []struct {
		Left  []A
		Right []A
	}{}
	for {
		leftCombo := []A{}
		rightCombo := []A{}
		for i := range xs {
			if current.Get(i) {
				leftCombo = append(leftCombo, xs[i])
			} else {
				rightCombo = append(rightCombo, xs[i])
			}
		}
		combo := struct {
			Left  []A
			Right []A
		}{
			Left:  leftCombo,
			Right: rightCombo,
		}
		combos = append(combos, combo)
		if current.Equal(max) {
			break
		}
		current = current.Inc()
	}
	return combos
}

func Subsets[A any](xs []A) [][]A {
	size := len(xs)
	if size == 0 {
		return nil
	}
	max := sets.NewBits(size)
	for i := 0; i < size; i++ {
		max.Set(i, true)
	}
	current := sets.NewBits(len(xs))
	// do not include empty set
	current = current.Inc()

	combos := [][]A{}
	for {
		combo := []A{}
		for i := range xs {
			if current.Get(i) {
				combo = append(combo, xs[i])
			}
		}
		combos = append(combos, combo)

		if current.Equal(max) {
			break
		}
		current = current.Inc()
	}
	return combos
}
