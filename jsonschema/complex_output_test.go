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
)

const SchemaComplexOutput = `
  {
    "$anchor": "output",
    "type": "object",
    "definitions": {
      "base58": {
        "$anchor": "base58",
        "type": "string",
        "pattern": "^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]+$"
      },
      "hex": {
        "$anchor": "hex",
        "type": "string",
        "pattern": "^[0123456789A-Fa-f]+$"
      },
      "tx_id": {
        "$anchor": "tx_id",
        "allOf": [
          {"$ref": "#hex"},
          {
            "minLength": 64,
            "maxLength": 64
          }
        ]
      },
      "address": {
        "$anchor": "address",
        "allOf": [
          {"$ref": "#base58"},
          {
            "minLength": 34,
            "maxLength": 34
          }
        ]
      },
      "signature": {
        "$anchor": "signature",
        "allOf": [
          {"$ref": "#hex"},
          {
            "minLength": 128,
            "maxLength": 128
          }
        ]
      }
    },
    "additionalProperties": false,
    "required": ["hash", "index", "value", "script"],
    "properties": {
      "hash": {"$ref": "#tx_id"},
      "index": {
        "type": "integer",
        "minimum": 0
      },
      "value": {
        "type": "integer"
      },
      "script": {
        "type": "object",
        "properties": {
          "type": {
            "type": "string",
            "enum": ["standard", "p2sh"]
          },
          "asm": {
            "type": "string"
          }
        }
      },
      "address": {"$ref": "#address"},
      "metadata": {
        "type": "object",
        "dependencies": {
          "wallet_path": ["public_seeds"]
        },
        "properties": {
          "wallet_path": {
            "type": "string"
          },
          "public_seeds": {
            "type": "object",
            "minProperties": 1,
            "maxProperties": 3,
            "additionalProperties": {
              "anyOf": [{"$ref": "#base58"}, {"$ref": "#hex"}]
            }
          }
        }
      }
    }
  }`

func TestComplexOutputJSONMatchParser(t *testing.T) {
	sch := SchemaComplexOutput

	g, err := newGrammar([]byte(sch))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("translated to: %v", g.String())

}
