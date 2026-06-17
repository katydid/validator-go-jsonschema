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

import "testing"

func TestMultipleOf2(t *testing.T) {
	var d float64 = 2
	valid := map[string]float64{
		"int by int": 10,
	}
	invalid := map[string]float64{
		"int by int fail": 7,
	}
	for desc, input := range valid {
		if !isMultipleOf(input, d) {
			t.Fatalf("wanted multipleof true for %v: %s", input, desc)
		}
	}
	for desc, input := range invalid {
		if isMultipleOf(input, d) {
			t.Fatalf("wanted multipleof false for %v: %s", input, desc)
		}
	}
}

func TestMultipleOf1p5(t *testing.T) {
	d := 1.5
	valid := map[string]float64{
		"zero is multiple of anything": 0,
		"4.5 is multiple of 1.5":       4.5,
	}
	invalid := map[string]float64{
		"35 is not multiple of 1.5": 35,
	}
	for desc, input := range valid {
		if !isMultipleOf(input, d) {
			t.Fatalf("wanted multipleof true for %v: %s", input, desc)
		}
	}
	for desc, input := range invalid {
		if isMultipleOf(input, d) {
			t.Fatalf("wanted multipleof false for %v: %s", input, desc)
		}
	}
}

func TestMultipleOfSmallNumber(t *testing.T) {
	d := 0.0001
	valid := map[string]float64{
		"0.0075 is multiple of 0.0001": 0.0075,
	}
	invalid := map[string]float64{
		"0.00751 is not multiple of 0.0001": 0.00751,
	}
	for desc, input := range valid {
		if !isMultipleOf(input, d) {
			t.Fatalf("wanted multipleof true for %v: %s", input, desc)
		}
	}
	for desc, input := range invalid {
		if isMultipleOf(input, d) {
			t.Fatalf("wanted multipleof false for %v: %s", input, desc)
		}
	}
}

func TestMultipleOfTooLarge(t *testing.T) {
	if isMultipleOf(1e308, 0.123456789) {
		t.Fatal("expected not multiple of")
	}
}

func TestMultipleOfOverflow(t *testing.T) {
	if !isMultipleOf(1e308, 0.5) {
		t.Fatal("expected multiple of")
	}
}

func TestMultipleOfSmallMultipleOfLarge(t *testing.T) {
	if !isMultipleOf(12391239123, 1e-8) {
		t.Fatal("expected multiple of")
	}
}
