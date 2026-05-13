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

import (
	"testing"

	"github.com/katydid/parser-go-json/json"
	"github.com/katydid/parser-go/parse"
)

// https://json-schema.org/learn/json-schema-examples#address
// A schema representing an address, with optional properties for different address components
// which enforces that locality, region, and countryName are required,
// and if postOfficeBox or extendedAddress is provided, streetAddress must also be provided.
const SchemaJSONSchemaExampleAddress = `{
  "type": "object",
  "properties": {
    "PostOfficeBox": {
      "type": "string"
    },
    "ExtendedAddress": {
      "type": "string"
    },
    "StreetAddress": {
      "type": "string"
    },
    "Locality": {
      "type": "string"
    },
    "Region": {
      "type": "string"
    },
    "PostalCode": {
      "type": "string"
    },
    "CountryName": {
      "type": "string"
    }
  },
  "required": [ "Locality", "Region", "CountryName" ],
  "dependentRequired": {
    "PostOfficeBox": [ "StreetAddress" ],
    "ExtendedAddress": [ "StreetAddress" ]
  }
}`

var addressFails = []string{
	`{"PostOfficeBox":"P","ExtendedAddress":"","Locality":"2A6i59","Region":"bf","CountryName":"aQ","Other":"j"}`,                                             // StreetAddress missing, but it is dependentRequired
	`{"PostOfficeBox":"","ExtendedAddress":"s09okTtQE","Locality":"Wr","Region":"A00mf66p","PostalCode":"y06dclIm","CountryName":"qld539n","Other":"W73k4i"}`, // StreetAddress missing dependentRequired
}

func TestAddressJSONMatchParser(t *testing.T) {
	sch := SchemaJSONSchemaExampleAddress
	fails := addressFails
	var p parse.ParserWithInit = json.NewJSONSchemaParser()

	for _, input := range fails {
		p.Init([]byte(input))
		m, err := MatchParser([]byte(sch), p)
		if err != nil {
			t.Fatal(err)
		}
		if m {
			t.Errorf("expected false, but got match for %s", input)
		}
	}
}

func TestAddressReflectMatchParser(t *testing.T) {
	sch := SchemaJSONSchemaExampleAddress
	fails := addressFails
	var p parse.ParserWithInit = newReflectParser()

	for _, input := range fails {
		p.Init([]byte(input))
		m, err := MatchParser([]byte(sch), p)
		if err != nil {
			t.Fatal(err)
		}
		if m {
			t.Errorf("expected false, but got match for %s", input)
		}
	}
}

// (
//   tag(object):*
//   &tag(object):
//     (
//       (.Locality:*&.Region:*&.CountryName:*)
//       &(
//         ((.ExtendedAddress:*&.StreetAddress:*)|(!(ExtendedAddress):*)*)
//         &((.PostOfficeBox:*&.StreetAddress:*)|(!(PostOfficeBox):*)*)
//       )
//       &{
//         (
//           (CountryName:->stringType()|ExtendedAddress:->stringType()|Locality:->stringType()|PostOfficeBox:->stringType()|PostalCode:->stringType()|Region:->stringType()|StreetAddress:->stringType())
//         )*
//         ;(
//           !((CountryName|ExtendedAddress|Locality|PostOfficeBox|PostalCode|Region|StreetAddress)):*
//         )*
//       }
//   )
// )

func TestAddressJSONCompile(t *testing.T) {
	sch := SchemaJSONSchemaExampleAddress
	fails := addressFails
	var p parse.ParserWithInit = json.NewJSONSchemaParser()

	fails = []string{
		`{"PostOfficeBox":"P","ExtendedAddress":"a","Locality":"2A6i59","Region":"bf","CountryName":"aQ","Other":"j"}`,
	}

	g, err := newGrammar([]byte(sch))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("translated to: %v", g.String())

	matcher, err := Compile([]byte(sch))
	if err != nil {
		t.Fatal(err)
	}

	for _, input := range fails {
		p.Init([]byte(input))
		m, err := matcher.MatchParser(p)
		if err != nil {
			t.Fatal(err)
		}
		if m {
			t.Errorf("expected false, but got match for %s", input)
		}
	}
}

func TestSmallAddressJSONCompile(t *testing.T) {
	sch := `{
  "type": "object",
  "properties": {
	"CountryName": {
      "type": "string"
    },
	"StreetAddress": {
      "type": "string"
    },
    "ExtendedAddress": {
      "type": "string"
    }
  },
  "dependentRequired": {
    "ExtendedAddress": [ "StreetAddress" ]
  }
}`

	fails := addressFails
	var p parse.ParserWithInit = json.NewJSONSchemaParser()

	fails = []string{
		`{"ExtendedAddress":"a"}`,
	}

	// 	translated := `
	// (
	// 	tag(object):*
	// 	&tag(object):
	// 		(
	// 			((.ExtendedAddress:*&.StreetAddress:*)|(!(ExtendedAddress):*)*)
	// 			&{
	// 				((CountryName:->stringType()|ExtendedAddress:->stringType()|StreetAddress:->stringType()))*
	// 				;(!((CountryName|ExtendedAddress|StreetAddress)):*)*
	// 			}
	// 		)
	// )`

	// 		(
	// 			.StreetAddress:*
	// 			&{
	// 				((CountryName:->stringType()|ExtendedAddress:->stringType()|StreetAddress:->stringType()))*
	// 				;(!((CountryName|ExtendedAddress|StreetAddress)):*)*
	// 			}
	// 		)

	if m, err := MatchBytes([]byte(sch), []byte(fails[0])); err != nil {
		panic(err)
	} else {
		if m {
			panic("didnt expect match")
		}
	}

	g, err := newGrammar([]byte(sch))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("translated to: %v", g.String())

	matcher, err := Compile([]byte(sch))
	if err != nil {
		t.Fatal(err)
	}

	for _, input := range fails {
		p.Init([]byte(input))
		m, err := matcher.MatchParser(p)
		if err != nil {
			t.Fatal(err)
		}
		if m {
			t.Errorf("expected false, but got match for %s", input)
		}
	}
}
