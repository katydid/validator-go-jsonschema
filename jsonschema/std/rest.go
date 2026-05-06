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

// Rest returns all elements in a slice, except for the index
func Rest[A any](xs []A, index int) []A {
	ys := make([]A, index)
	copy(ys, xs)
	return append(ys, xs[index+1:]...)
}

// Rests combinations of slices each with one element missing.
func Rests[A any](xs []A) [][]A {
	rs := make([][]A, len(xs))
	for i := range xs {
		rs[i] = Rest(xs, i)
	}
	return rs
}
