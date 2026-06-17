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

import (
	"testing"

	"github.com/katydid/validator-go-jsonschema/jsonschema/funcs/regexformat"
)

type regexSuite struct {
	description string
	schema      regexSchema
	tests       []test
}

type regexSchema struct {
	typ     string
	pattern string
}

type test struct {
	description string
	data        string
	valid       bool
}

func TestRegexECMASuites(t *testing.T) {
	for _, suite := range emcaSuites {
		matchString, err := regexformat.Compile(suite.schema.pattern)
		if err != nil {
			t.Errorf("error compiling regex for: %s: %s: %v", suite.description, suite.schema.pattern, err)
		} else {
			for _, test := range suite.tests {
				got := matchString(test.data)
				if got != test.valid {
					t.Errorf("%s: %s: %s: %s want %v got %v", suite.description, suite.schema.pattern, test.description, test.data, test.valid, got)
				}
			}
		}
	}
}

// Copied from draft4/optional/ecmascript-regex.json
var emcaSuites []regexSuite = []regexSuite{
	{
		description: "ECMA 262 regex $ does not match trailing newline",
		schema: regexSchema{
			typ:     "string",
			pattern: "^abc$",
		},
		tests: []test{
			{
				description: "matches in Python, but not in ECMA 262",
				data:        "abc\\n",
				valid:       false,
			},
			{
				description: "matches",
				data:        "abc",
				valid:       true,
			},
		},
	},
	{
		description: "ECMA 262 regex converts \\t to horizontal tab",
		schema: regexSchema{
			typ:     "string",
			pattern: "^\\t$",
		},
		tests: []test{
			{
				description: "does not match",
				data:        "\\t",
				valid:       false,
			},
			{
				description: "matches",
				data:        "\u0009",
				valid:       true,
			},
		},
	},
	{
		description: "ECMA 262 regex escapes control codes with \\c and upper letter",
		schema: regexSchema{
			typ:     "string",
			pattern: "^\\cC$",
		},
		tests: []test{
			{
				description: "does not match",
				data:        "\\cC",
				valid:       false,
			},
			{
				description: "matches",
				data:        "\u0003",
				valid:       true,
			},
		},
	},
	{
		description: "ECMA 262 regex escapes control codes with \\c and lower letter",
		schema: regexSchema{
			typ:     "string",
			pattern: "^\\cc$",
		},
		tests: []test{
			{
				description: "does not match",
				data:        "\\cc",
				valid:       false,
			},
			{
				description: "matches",
				data:        "\u0003",
				valid:       true,
			},
		},
	},
	{
		description: "ECMA 262 \\d matches ascii digits only",
		schema: regexSchema{
			typ:     "string",
			pattern: "^\\d$",
		},
		tests: []test{
			{
				description: "ASCII zero matches",
				data:        "0",
				valid:       true,
			},
			{
				description: "NKO DIGIT ZERO does not match (unlike e.g. Python)",
				data:        "߀",
				valid:       false,
			},
			{
				description: "NKO DIGIT ZERO (as \\u escape) does not match",
				data:        "\u07c0",
				valid:       false,
			},
		},
	},
	{
		description: "ECMA 262 \\D matches everything but ascii digits",
		schema: regexSchema{
			typ:     "string",
			pattern: "^\\D$",
		},
		tests: []test{
			{
				description: "ASCII zero does not match",
				data:        "0",
				valid:       false,
			},
			{
				description: "NKO DIGIT ZERO matches (unlike e.g. Python)",
				data:        "߀",
				valid:       true,
			},
			{
				description: "NKO DIGIT ZERO (as \\u escape) matches",
				data:        "\u07c0",
				valid:       true,
			},
		},
	},
	{
		description: "ECMA 262 \\w matches ascii letters only",
		schema: regexSchema{
			typ:     "string",
			pattern: "^\\w$",
		},
		tests: []test{
			{
				description: "ASCII 'a' matches",
				data:        "a",
				valid:       true,
			},
			{
				description: "latin-1 e-acute does not match (unlike e.g. Python)",
				data:        "é",
				valid:       false,
			},
		},
	},
	{
		description: "ECMA 262 \\W matches everything but ascii letters",
		schema: regexSchema{
			typ:     "string",
			pattern: "^\\W$",
		},
		tests: []test{
			{
				description: "ASCII 'a' does not match",
				data:        "a",
				valid:       false,
			},
			{
				description: "latin-1 e-acute matches (unlike e.g. Python)",
				data:        "é",
				valid:       true,
			},
		},
	},
	{
		description: "ECMA 262 \\s matches whitespace",
		schema: regexSchema{
			typ:     "string",
			pattern: "^\\s$",
		},
		tests: []test{
			{
				description: "ASCII space matches",
				data:        " ",
				valid:       true,
			},
			{
				description: "Character tabulation matches",
				data:        "\t",
				valid:       true,
			},
			{
				description: "Line tabulation matches",
				data:        "\u000b",
				valid:       true,
			},
			{
				description: "Form feed matches",
				data:        "\u000c",
				valid:       true,
			},
			{
				description: "latin-1 non-breaking-space matches",
				data:        "\u00a0",
				valid:       true,
			},
			{
				description: "zero-width whitespace matches",
				data:        "\ufeff",
				valid:       true,
			},
			{
				description: "line feed matches (line terminator)",
				data:        "\u000a",
				valid:       true,
			},
			{
				description: "paragraph separator matches (line terminator)",
				data:        "\u2029",
				valid:       true,
			},
			{
				description: "EM SPACE matches (Space_Separator)",
				data:        "\u2003",
				valid:       true,
			},
			{
				description: "Non-whitespace control does not match",
				data:        "\u0001",
				valid:       false,
			},
			{
				description: "Non-whitespace does not match",
				data:        "\u2013",
				valid:       false,
			},
		},
	},
	{
		description: "ECMA 262 \\S matches everything but whitespace",
		schema: regexSchema{
			typ:     "string",
			pattern: "^\\S$",
		},
		tests: []test{
			{
				description: "ASCII space does not match",
				data:        " ",
				valid:       false,
			},
			{
				description: "Character tabulation does not match",
				data:        "\t",
				valid:       false,
			},
			{
				description: "Line tabulation does not match",
				data:        "\u000b",
				valid:       false,
			},
			{
				description: "Form feed does not match",
				data:        "\u000c",
				valid:       false,
			},
			{
				description: "latin-1 non-breaking-space does not match",
				data:        "\u00a0",
				valid:       false,
			},
			{
				description: "zero-width whitespace does not match",
				data:        "\ufeff",
				valid:       false,
			},
			{
				description: "line feed does not match (line terminator)",
				data:        "\u000a",
				valid:       false,
			},
			{
				description: "paragraph separator does not match (line terminator)",
				data:        "\u2029",
				valid:       false,
			},
			{
				description: "EM SPACE does not match (Space_Separator)",
				data:        "\u2003",
				valid:       false,
			},
			{
				description: "Non-whitespace control matches",
				data:        "\u0001",
				valid:       true,
			},
			{
				description: "Non-whitespace matches",
				data:        "\u2013",
				valid:       true,
			},
		},
	},
	{
		description: "patterns always use unicode semantics with pattern",
		schema: regexSchema{
			pattern: "\\p{Letter}cole",
		},
		tests: []test{
			{
				description: "ascii character in json string",
				data:        "Les hivers de mon enfance etaient des saisons longues, longues. Nous vivions en trois lieux: l'ecole, l'eglise et la patinoire; mais la vraie vie etait sur la patinoire.",
				valid:       true,
			},
			{
				description: "literal unicode character in json string",
				data:        "Les hivers de mon enfance étaient des saisons longues, longues. Nous vivions en trois lieux: l'école, l'église et la patinoire; mais la vraie vie était sur la patinoire.",
				valid:       true,
			},
			{
				description: "unicode character in hex format in string",
				data:        "Les hivers de mon enfance étaient des saisons longues, longues. Nous vivions en trois lieux: l'\u00e9cole, l'église et la patinoire; mais la vraie vie était sur la patinoire.",
				valid:       true,
			},
			{
				description: "unicode matching is case-sensitive",
				data:        "LES HIVERS DE MON ENFANCE ÉTAIENT DES SAISONS LONGUES, LONGUES. NOUS VIVIONS EN TROIS LIEUX: L'ÉCOLE, L'ÉGLISE ET LA PATINOIRE; MAIS LA VRAIE VIE ÉTAIT SUR LA PATINOIRE.",
				valid:       false,
			},
		},
	},
	{
		description: "\\w in patterns matches [A-Za-z0-9_], not unicode letters",
		schema: regexSchema{
			pattern: "\\wcole",
		},
		tests: []test{
			{
				description: "ascii character in json string",
				data:        "Les hivers de mon enfance etaient des saisons longues, longues. Nous vivions en trois lieux: l'ecole, l'eglise et la patinoire; mais la vraie vie etait sur la patinoire.",
				valid:       true,
			},
			{
				description: "literal unicode character in json string",
				data:        "Les hivers de mon enfance étaient des saisons longues, longues. Nous vivions en trois lieux: l'école, l'église et la patinoire; mais la vraie vie était sur la patinoire.",
				valid:       false,
			},
			{
				description: "unicode character in hex format in string",
				data:        "Les hivers de mon enfance étaient des saisons longues, longues. Nous vivions en trois lieux: l'\u00e9cole, l'église et la patinoire; mais la vraie vie était sur la patinoire.",
				valid:       false,
			},
			{
				description: "unicode matching is case-sensitive",
				data:        "LES HIVERS DE MON ENFANCE ÉTAIENT DES SAISONS LONGUES, LONGUES. NOUS VIVIONS EN TROIS LIEUX: L'ÉCOLE, L'ÉGLISE ET LA PATINOIRE; MAIS LA VRAIE VIE ÉTAIT SUR LA PATINOIRE.",
				valid:       false,
			},
		},
	},
	{
		description: "pattern with ASCII ranges",
		schema: regexSchema{
			pattern: "[a-z]cole",
		},
		tests: []test{
			{
				description: "literal unicode character in json string",
				data:        "Les hivers de mon enfance étaient des saisons longues, longues. Nous vivions en trois lieux: l'école, l'église et la patinoire; mais la vraie vie était sur la patinoire.",
				valid:       false,
			},
			{
				description: "unicode character in hex format in string",
				data:        "Les hivers de mon enfance étaient des saisons longues, longues. Nous vivions en trois lieux: l'\u00e9cole, l'église et la patinoire; mais la vraie vie était sur la patinoire.",
				valid:       false,
			},
			{
				description: "ascii characters match",
				data:        "Les hivers de mon enfance etaient des saisons longues, longues. Nous vivions en trois lieux: l'ecole, l'eglise et la patinoire; mais la vraie vie etait sur la patinoire.",
				valid:       true,
			},
		},
	},
	{
		description: "\\d in pattern matches [0-9], not unicode digits",
		schema: regexSchema{
			pattern: "^\\d+$",
		},
		tests: []test{
			{
				description: "ascii digits",
				data:        "42",
				valid:       true,
			},
			{
				description: "ascii non-digits",
				data:        "-%#",
				valid:       false,
			},
			{
				description: "non-ascii digits (BENGALI DIGIT FOUR, BENGALI DIGIT TWO)",
				data:        "৪২",
				valid:       false,
			},
		},
	},
	{
		description: "pattern with non-ASCII digits",
		schema: regexSchema{
			pattern: "^\\p{digit}+$",
		},
		tests: []test{
			{
				description: "ascii digits",
				data:        "42",
				valid:       true,
			},
			{
				description: "ascii non-digits",
				data:        "-%#",
				valid:       false,
			},
			{
				description: "non-ascii digits (BENGALI DIGIT FOUR, BENGALI DIGIT TWO)",
				data:        "৪২",
				valid:       true,
			},
		},
	},
}

// draft4/optional/emcascript-regex.json
// ECMA 262 regex $ does not match trailing newline
func TestRegexECMATrailingNewline(t *testing.T) {
	matchString, err := regexformat.Compile("^abc$")
	if err != nil {
		t.Fatal(err)
	}
	valid := map[string]string{
		"matches": "abc",
	}
	invalid := map[string]string{
		"matches in Python, but not in ECMA 262": "abc\\n",
	}
	for desc, input := range valid {
		m := matchString(input)
		if !m {
			t.Fatalf("expected match for %s", desc)
		}
	}
	for desc, input := range invalid {
		m := matchString(input)
		if m {
			t.Fatalf("expected no match for %s", desc)
		}
	}
}

// Copied from draft4/optional/non-bmp-regex.json
// Proper UTF-16 surrogate pair handling: pattern
// Optional because .Net doesn't correctly handle 32-bit Unicode characters
func TestRegexUTF16(t *testing.T) {
	matchString, err := regexformat.Compile("^🐲*$")
	if err != nil {
		t.Fatal(err)
	}
	valid := map[string]string{
		"matches empty":  "",
		"matches single": "🐲",
		"matches two":    "🐲🐲",
	}
	invalid := map[string]string{
		"doesn't match one":       "🐉",
		"doesn't match two":       "🐉🐉",
		"doesn't match one ASCII": "D",
		"doesn't match two ASCII": "DD",
	}
	for desc, input := range valid {
		m := matchString(input)
		if !m {
			t.Fatalf("expected match for %s", desc)
		}
	}
	for desc, input := range invalid {
		m := matchString(input)
		if err != nil {
			t.Fatal(err)
		}
		if m {
			t.Fatalf("expected no match for %s", desc)
		}
	}
}
