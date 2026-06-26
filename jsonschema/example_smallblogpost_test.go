package jsonschema

import (
	"testing"

	"github.com/katydid/parser-go-json/json"
	"github.com/katydid/parser-go/parse"
)

const SchemaSmallBlogPostExample = `{ "title": "small jsonschema for a blogpost",
  "type":"object", "additionalProperties":false, "required": ["content"],
  "properties": {
    "content": { "type":"string" },
    "author": { "$ref":"#/definitions/user-profile" } },
  "definitions": { "user-profile": {
    "type": "object", "additionalProperties":false, "required": ["username"], 
    "properties": {
      "username": { "type":"string" },
      "email": { "type":"string", "format":"email" } } } } }`

func TestSchemaSmallBlogPostExample(t *testing.T) {
	sch := SchemaSmallBlogPostExample
	passes := []string{
		`{"content": "Dragons"}`,
		`{"content": "Dragons", "author": {"username": "Khaleesi"}}`,
	}
	var p parse.ParserWithInit = json.NewJSONSchemaParser()

	g, err := newGrammar([]byte(sch))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("translated to: %v", g.String())
	for _, input := range passes {
		p.Init([]byte(input))
		m, err := MatchParser([]byte(sch), p)
		if err != nil {
			t.Fatal(err)
		}
		if !m {
			t.Errorf("expected true, but got no match for %s", input)
		}
	}
}
