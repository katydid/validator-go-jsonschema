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
	"strings"
	"testing"

	"github.com/katydid/parser-go-json/json"
)

func BenchmarkComplexValid(b *testing.B) {
	name := "jsck-complex-valid"
	want := !strings.Contains(name, "-invalid")
	benchmarks, err := getBenchmarks()
	if err != nil {
		b.Fatal(err)
	}
	ran := false
	for _, suite := range benchmarks {
		if suite.name != name {
			continue
		}
		matcher, err := Compile(suite.schema)
		if err != nil {
			b.Fatal(err)
		}
		parser := json.NewJSONSchemaParser()
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			ran = true
			for i := 0; i < b.N; i++ {
				data := suite.datas[i%len(suite.datas)]
				parser.Init(data)
				got, err := matcher.MatchParser(parser)
				if err != nil {
					b.Fatal(err)
				}
				if got != want {
					b.Errorf("want %v, but got %v for instance: %s", want, got, data)
				}
			}
			b.ReportAllocs()
		})
		if !ran {
			b.Fatalf("did not find %s in %v", name, benchmarks)
		}
	}
}
