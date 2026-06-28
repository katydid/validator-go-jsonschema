# Json Schema using Brzozowski's derivatives

JSON Schema in Go using Katydid underlying algorithm.
This project translates JSON Schema to a regular hedge grammar and then executes validation via Katydid's [validator-go](github.com/katydid/validator-go).

## Usage Example

```go
	schemaBytes := []byte(`
    { "title": "small jsonschema for a blogpost",
      "type":"object", "additionalProperties":false, "required": ["content"],
      "properties": {
        "content": { "type":"string" },
        "author": { "$ref":"#/definitions/user-profile" } },
      "definitions": { "user-profile": {
        "type": "object", "additionalProperties":false, "required": ["username"], 
        "properties": {
          "username": { "type":"string" },
          "email": { "type":"string", "format":"email" } } } } }`)
	compiled, err := Compile(schemaBytes)
	...
	input := []byte(`{"content": "Dragons", "author": {"username": "Khaleesi"}}`)
	matched, err := compiled.MatchBytes(input)
	...
```

## Test Suites passed

* Draft4 (excluding `uniqueItems` and `remoteRef`)

## Unsupported

* [uniqueItems](./decisions/uniqueItems.md)