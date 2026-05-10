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

package jsonschema

import "testing"

const SchemaConfIsIn2026OrLate2025AndEU = `{
		"definitions": {
			"due": {
				"type": "object",
				"anyOf": [
				{
					"properties": {
						"Year": {
							"$ref": "#/definitions/year2026"
						}
					},
					"required": ["Year"]
				},
				{
					"allOf": [                  
						{
							"properties": {
								"Year": {
									"$ref": "#/definitions/year2025"
								}
							},
							"required": ["Year"]
						},
						{
							"properties": {
								"Month": {
									"$ref": "#/definitions/month10"
								}
							},
							"required": ["Month"]
						}
					]
				}
				]
			},
			"loc": {
				"type": "object",
				"properties": {
					"Cont": {
					"$ref": "#/definitions/conteu"
						}
					},
				"required": ["Cont"]
			},
			"year2026": {
				"const": "2026"
			},
			"year2025": {
				"const": "2025"
			},
			"month10": {
				"minimum": 10
			},
			"conteu": {
				"const": "EU"
			}
		},
		"type": "object",
		"properties": {
			"Due": {
				"$ref": "#/definitions/due"
			},
			"Loc": {
				"$ref": "#/definitions/loc"
			}
		},
		"required": ["Due", "Loc"]
	}`

var confFails = []string{
	`{"Name":"Ich","Due":{"Year":"2052","Month":"04","Day":"19"},"Notify":{"Year":"JPPsvD","Month":null,"Day":"j"},"Loc":{"Cont":"RzgAnuS","Ctry":"2697S","City":"bQC"},"Category":"f7J3Pb","Tags":["L64","QVMoWrb6l"]}`,
}

func TestConf(t *testing.T) {
	g, err := newGrammar([]byte(SchemaConfIsIn2026OrLate2025AndEU))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("translated to: %v", g.String())
	for _, input := range confFails {
		m, err := MatchBytes([]byte(SchemaConfIsIn2026OrLate2025AndEU), []byte(input))
		if err != nil {
			t.Fatal(err)
		}
		if m {
			t.Errorf("expected false, but got match for %s", input)
		}
	}
}
