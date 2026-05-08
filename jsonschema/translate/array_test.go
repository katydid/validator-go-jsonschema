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

package translate

import (
	"reflect"
	"testing"

	"github.com/katydid/validator-go/validator/ast"
)

func TestConcatCombos(t *testing.T) {
	expect := func(input []*ast.Pattern, want []*ast.Pattern) {
		t.Helper()
		got := concatCombos(input)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v got %v", want, got)
		}
	}
	z := ast.NewZAny()
	expect([]*ast.Pattern{}, []*ast.Pattern{ast.NewEmpty()})
	expect([]*ast.Pattern{ast.NewNot(z)}, []*ast.Pattern{ast.NewEmpty(), ast.NewConcat(ast.NewNot(z), ast.NewZAny())})
	expect([]*ast.Pattern{ast.NewNot(z), ast.NewContains(z)}, []*ast.Pattern{ast.NewEmpty(), ast.NewNot(z), ast.NewConcat(ast.NewNot(z), ast.NewContains(z), ast.NewZAny())})
	expect(
		[]*ast.Pattern{ast.NewNot(z), ast.NewContains(z), ast.NewOptional(z)},
		[]*ast.Pattern{ast.NewEmpty(), ast.NewNot(z), ast.NewConcat(ast.NewNot(z), ast.NewContains(z)), ast.NewConcat(ast.NewNot(z), ast.NewContains(z), ast.NewOptional(z), ast.NewZAny())},
	)
}
