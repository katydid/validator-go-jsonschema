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

func TestRest(t *testing.T) {
	testutil.ExpectErr(t, "", Rest([]int{1, 2, 3}, 1), []int{1, 3})
	testutil.ExpectErr(t, "", Rest([]int{1, 2, 3, 4}, 1), []int{1, 3, 4})
	testutil.ExpectErr(t, "", Rest([]int{1, 2, 3, 4}, 0), []int{2, 3, 4})
	testutil.ExpectErr(t, "", Rest([]int{1, 2, 3, 4}, 3), []int{1, 2, 3})
}

func TestRests(t *testing.T) {
	testutil.ExpectErr(t, "", Rests([]int{1, 2, 3}), [][]int{{2, 3}, {1, 3}, {1, 2}})
}
