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

package datetimeformat

import "testing"

func TestDateTimeFormat(t *testing.T) {
	tests := map[string]bool{
		"1963-06-19T08:30:06.283185Z":   true,
		"1963-06-19T08:30:06Z":          true,
		"1937-01-01T12:00:27.87+00:20":  true,
		"1990-12-31T15:59:50.123-08:00": true,
		"1998-12-31T23:59:60Z":          true,
		"1998-12-31T15:59:60.123-08:00": true,
		"1963-06-19t08:30:06.283185z":   true,
	}
	for test, valid := range tests {
		t.Run(test, func(t *testing.T) {
			if IsValid([]byte(test)) != valid {
				t.Error()
			}
		})
	}
}
