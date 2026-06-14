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

package hostname

import (
	"testing"

	"github.com/katydid/validator-go-jsonschema/jsonschema/funcs/hostname/lexer"
)

var testData = map[string]bool{
	"en.wikipedia.org":                true,
	"www.example.com":                 true,
	"xn--4gbwdl.xn--wgbh1c":           true,
	"ab--cd.example":                  true,
	"-a-host-name-that-starts-with--": false,
	"not_a_valid_host_name":           false,
	"-hostname":                       false,
	"hostname-":                       false,
	"_hostname":                       false,
	"hostname_":                       false,
	"host_name":                       false,
	"a.b":                             true,
	"a-b.cde":                         true,
	"127.0.0.1":                       true,
	"___":                             false,
}

func TestIsValid(t *testing.T) {
	l := lexer.NewLexer([]byte(nil))
	for input, ok := range testData {
		if ok {
			if !l.IsValid([]byte(input)) {
				t.Errorf("want valid %s", input)
			}
		} else {
			if l.IsValid([]byte(input)) {
				t.Errorf("want not valid %s", input)
			}
		}
	}
}
