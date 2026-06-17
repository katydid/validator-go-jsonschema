//  Copyright 2015 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package jsonschema

import (
	"testing"
)

const pathDraft4 = "../../../json-schema-org/JSON-Schema-Test-Suite/tests/draft4/"

var supportedDraft4 = &Supported{
	passingFiles: map[string]bool{
		"additionalItems.json":      true,
		"additionalProperties.json": true,
		"allOf.json":                true,
		"anyOf.json":                true,
		"default.json":              true,
		// "dependencies.json": true,
		"enum.json":                    true,
		"format.json":                  true,
		"infinite-loop-detection.json": true,
		"items.json":                   true,
		"maximum.json":                 true,
		"maxItems.json":                true,
		"maxLength.json":               true,
		"maxProperties.json":           true,
		"minimum.json":                 true,
		"minItems.json":                true,
		"minLength.json":               true,
		"minProperties.json":           true,
		"multipleOf.json":              true,
		"not.json":                     true,
		"oneOf.json":                   true,
		"pattern.json":                 true,
		"patternProperties.json":       true,
		"properties.json":              true,
		// "ref.json": true,
		"required.json": true,
		"type.json":     true,

		// optional
		"optional/ecmascript-regex.json":     true,
		"optional/non-bmp-regex.json":        true,
		"optional/zeroTerminatedFloats.json": true,

		// optional/format
		"optional/bignum.json":           true,
		"optional/format/date-time.json": true,
		"optional/format/email.json":     true,
		"optional/format/hostname.json":  true,
		"optional/format/ipv4.json":      true,
		"optional/format/ipv6.json":      true,
		"optional/format/unknown.json":   true,
		"optional/format/uri.json":       true,
	},
	skippingFiles: map[string]bool{
		"uniqueItems.json": true, // We do not support uniqueItems, see https://github.com/katydid/validator-go-jsonschema/blob/main/decisions/uniqueItems.md
		"refRemote.json":   true, // remote and file ref support should be relatively easy to add, but is just not of theoretical importance at this stage.
		"definitions.json": true, // remote and file ref support should be relatively easy to add, but is just not of theoretical importance at this stage.
		// optional
		"optional/id.json":             true, // remote and file ref support should be relatively easy to add, but is just not of theoretical importance at this stage.
		"optional/float-overflow.json": true, // Need better checking for float overflow to convert to decimal in the json parser and we need to support decimal in multipleOf
	},
	passingTests:  map[string]bool{},
	skippingTests: map[string]bool{},
}

func TestSuiteDraft4(t *testing.T) {
	runTests(t, pathDraft4, supportedDraft4)
}
