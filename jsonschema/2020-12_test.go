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

const path202012 = "../../../json-schema-org/JSON-Schema-Test-Suite/tests/draft2020-12/"

var supported202012 = &Supported{
	passingFiles: map[string]bool{
		// "additionalProperties.json":    true,
		// "allOf.json": true,
		// "anchor.json":                  true,
		// "anyOf.json":                   true,
		// "boolean_schema.json":          true,
		// "const.json":                   true,
		// "contains.json":                true,
		// "content.json":                 true,
		// "default.json":                 true,
		// "defs.json":                    true,
		// "dependentRequired.json":       true,
		// "dependentSchemas.json":        true,
		// "dynamicRef.json":              true,
		// "enum.json":                    true,
		// "exclusiveMaximum.json": true,
		// "exclusiveMinimum.json":        true,
		"format.json": true,
		// "if-then-else.json":            true,
		// "infinite-loop-detection.json": true,
		// "items.json":                   true,
		// "maxContains.json":             true,
		"maximum.json": true,
		// "maxItems.json":                true,
		// "maxLength.json":               true,
		// "maxProperties.json":           true,
		// "minContains.json":             true,
		"minimum.json": true,
		// "minItems.json":                true,
		// "minLength.json":               true,
		// "minProperties.json":           true,
		"multipleOf.json": true,
		// "not.json":                     true,
		// "oneOf.json":                   true,
		// "pattern.json":                 true,
		// "patternProperties.json":       true,
		// "prefixItems.json":             true,
		// "properties.json": true,
		// "propertyNames.json": true,
		// "ref.json":                     true,
		// "refRemote.json":               true,
		// "required.json":                true,
		// "type.json":                    true,
		// "unevaluatedItems.json":        true,
		// "unevaluatedProperties.json":   true,
		// "uniqueItems.json":             true,
		// "vocabulary.json":              true,
	},
	skippingFiles: map[string]bool{
		"dependencies-compatibility.json": true,
	},
	passingTests: map[string]bool{},
	skippingTests: map[string]bool{
		// TODO: not sure about this, but  skipping for now
		"format.json:ipv4 format:invalid ipv4 string is only an annotation by default":                                   true,
		"format.json:email format:invalid email string is only an annotation by default":                                 true,
		"format.json:ipv6 format:invalid ipv6 string is only an annotation by default":                                   true,
		"format.json:hostname format:invalid hostname string is only an annotation by default":                           true,
		"format.json:date format:invalid date string is only an annotation by default":                                   true,
		"format.json:date-time format:invalid date-time string is only an annotation by default":                         true,
		"format.json:time format:invalid time string is only an annotation by default":                                   true,
		"format.json:json-pointer format:invalid json-pointer string is only an annotation by default":                   true,
		"format.json:relative-json-pointer format:invalid relative-json-pointer string is only an annotation by default": true,
		"format.json:iri format:invalid iri string is only an annotation by default":                                     true,
		"format.json:iri-reference format:invalid iri-reference string is only an annotation by default":                 true,
		"format.json:uri format:invalid uri string is only an annotation by default":                                     true,
		"format.json:uri-reference format:invalid uri-reference string is only an annotation by default":                 true,
		"format.json:uri-template format:invalid uri-template string is only an annotation by default":                   true,
		"format.json:uuid format:invalid uuid string is only an annotation by default":                                   true,
		"format.json:duration format:invalid duration string is only an annotation by default":                           true,
	},
}

func TestSuite202012(t *testing.T) {
	runTests(t, path202012, supported202012)
}
