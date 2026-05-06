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

package translate

import (
	"fmt"
	"sort"

	"github.com/katydid/validator-go-jsonschema/jsonschema/schema"
	"github.com/katydid/validator-go-jsonschema/jsonschema/std"
	"github.com/katydid/validator-go/validator/ast"
	"github.com/katydid/validator-go/validator/combinator"
)

func TranslateDraft4(schema *schema.Schema) (*ast.Grammar, error) {
	p, err := translate(schema)
	if err != nil {
		return nil, err
	}
	return ast.NewGrammar(ast.RefLookup(map[string]*ast.Pattern{"main": p})), nil
}

func translate(schema *schema.Schema) (*ast.Pattern, error) {
	pattern, err := translateOne(schema)
	if err != nil {
		return nil, err
	}
	if schema.Type != nil {
		types := *schema.Type
		if len(types) == 1 {
			p, err := translateType(types[0])
			if err != nil {
				return nil, err
			}
			pattern = ast.NewAnd(p, pattern)
		} else {
			ps := make([]*ast.Pattern, len(types))
			for i := range types {
				var err error
				ps[i], err = translateType(types[i])
				if err != nil {
					return nil, err
				}
			}
			ors := ast.NewOr(ps...)
			pattern = ast.NewAnd(ors, pattern)
		}
	}
	return pattern, nil
}

func translateOne(schema *schema.Schema) (*ast.Pattern, error) {
	if len(schema.Id) > 0 {
		return nil, fmt.Errorf("id not supported")
	}
	if schema.Default != nil {
		return nil, fmt.Errorf("default not supported")
	}
	if schema.HasNumericConstraints() {
		p, err := translateNumeric(schema.Numeric)
		return p, err
	}
	if schema.HasStringConstraints() {
		if schema.Type != nil && len(*schema.Type) > 1 {
			return nil, fmt.Errorf("list of types not supported with string constraints %#v", schema)
		}
		p, err := translateString(schema.String, schema.Format)
		return p, err
	}
	if schema.HasArrayConstraints() {
		return nil, fmt.Errorf("array not supported")
	}
	if schema.HasObjectConstraints() {
		p, err := translateObject(schema)
		return p, err
	}
	if schema.HasOperatorConstraints() {
		p, err := translateInstance(schema)
		return p, err
	}
	if len(schema.Format) > 0 {
		expr, err := translateFormat(schema.Format)
		if err != nil {
			return nil, err
		}
		return combinator.Value(expr), nil
	}

	if len(schema.Ref) > 0 {
		return nil, fmt.Errorf("ref not supported")
	}
	return ast.NewZAny(), nil
}

func translates(schemas []*schema.Schema) ([]*ast.Pattern, error) {
	return std.MapErr(schemas, translate)
}

func translateInstance(schema *schema.Schema) (*ast.Pattern, error) {
	if len(schema.Definitions) > 0 {
		return nil, fmt.Errorf("definitions not supported")
	}
	if len(schema.Enum) > 0 {
		return nil, fmt.Errorf("enum not supported")
	}
	if len(schema.AllOf) > 0 {
		ps, err := translates(schema.AllOf)
		if err != nil {
			return nil, err
		}
		return ast.NewAnd(ps...), nil
	}
	if len(schema.AnyOf) > 0 {
		ps, err := translates(schema.AnyOf)
		if err != nil {
			return nil, err
		}
		return ast.NewOr(ps...), nil
	}
	if len(schema.OneOf) > 0 {
		return translateOneOf(schema.OneOf)
	}
	if schema.Not != nil {
		p, err := translate(schema.Not)
		if err != nil {
			return nil, err
		}
		return ast.NewNot(p), nil
	}
	panic("unreachable object")
}

func translateType(typ schema.SimpleType) (*ast.Pattern, error) {
	switch typ {
	case schema.TypeArray, schema.TypeObject:
		//This does not distinguish between arrays and objects
		return combinator.Many(combinator.InAny(combinator.Any())), nil
	case schema.TypeBoolean:
		return combinator.Value(boolTypeExpr()), nil
	case schema.TypeInteger:
		return combinator.Value(integerTypeExpr()), nil
	case schema.TypeNull:
		return combinator.Value(nullTypeExpr()), nil
	case schema.TypeNumber:
		return combinator.Value(numberTypeExpr()), nil
	case schema.TypeString:
		return combinator.Value(stringTypeExpr()), nil
	}
	panic(fmt.Sprintf("unknown simpletype: %s", typ))
}

func translateObject(s *schema.Schema) (*ast.Pattern, error) {
	if s.MaxProperties != nil {
		return nil, fmt.Errorf("maxProperties not supported")
	}
	if s.MinProperties > 0 {
		return nil, fmt.Errorf("minProperties not supported")
	}
	required := make(map[string]struct{})
	for _, req := range s.Required {
		required[req] = struct{}{}
	}
	requiredIf := make(map[string][]string)
	moreProperties := make(map[string]*schema.Schema)
	if s.Dependencies != nil {
		deps := *s.Dependencies
		for name, dep := range deps {
			if len(dep.RequiredProperty) > 0 {
				requiredIf[name] = deps[name].RequiredProperty
			} else {
				moreProperties[name] = deps[name].Schema
			}
		}
	}

	names := []string{}
	for name := range s.Properties {
		names = append(names, name)
	}
	sort.Strings(names)

	patternNames := []string{}
	for name := range s.PatternProperties {
		patternNames = append(patternNames, name)
	}
	sort.Strings(patternNames)

	additional := ast.NewZAny()
	if len(names) > 0 || len(patternNames) > 0 {
		nameExprs := make([]*ast.NameExpr, len(names)+len(patternNames))
		for i, name := range names {
			nameExprs[i] = ast.NewStringName(name)
		}
		for i, name := range patternNames {
			nameExprs[i+len(names)] = ast.NewRegexName(name)
		}
		additional = ast.NewZeroOrMore(
			ast.NewTreeNode(ast.NewAnyNameExcept(
				ast.NewNameChoice(nameExprs...),
			), ast.NewZAny()),
		)
	}
	if s.AdditionalProperties != nil {
		if s.AdditionalProperties.Bool != nil && !(*s.AdditionalProperties.Bool) {
			additional = ast.NewEmpty()
		} else if s.AdditionalProperties.Type != schema.TypeUnknown {
			typ, err := translateType(s.AdditionalProperties.Type)
			if err != nil {
				return nil, err
			}
			additional = ast.NewZeroOrMore(
				ast.NewTreeNode(ast.NewAnyName(), typ),
			)
		}
	}
	patterns := make(map[string]*ast.Pattern)
	for _, name := range names {
		child, err := translate(s.Properties[name])
		if err != nil {
			return nil, err
		}
		patterns[name] = ast.NewTreeNode(ast.NewStringName(name), child)
	}
	for _, name := range patternNames {
		child, err := translate(s.PatternProperties[name])
		if err != nil {
			return nil, err
		}
		patterns[name] = ast.NewTreeNode(ast.NewRegexName(name), child)
	}
	for _, name := range names {
		if _, ok := requiredIf[name]; ok {
			return nil, fmt.Errorf("dependencies are not supported")
		}
		if _, ok := moreProperties[name]; ok {
			return nil, fmt.Errorf("dependencies are not supported")
		}
		if _, ok := required[name]; !ok {
			patterns[name] = ast.NewOptional(patterns[name])
		}
	}

	patternList := make([]*ast.Pattern, 0, len(patterns))
	for _, name := range names {
		patternList = append(patternList, patterns[name])
	}
	patternList = append(patternList, additional)
	// TODO: Be more specific and create ast.NewTagName for "object"
	return ast.NewTreeNode(ast.NewStringName("object"), ast.NewInterleave(patternList...)), nil
}

func and(list []*ast.Expr) *ast.Expr {
	return std.MustFoldA(list, combinator.And)
}

func translateArray(s *schema.Schema) (*ast.Pattern, error) {
	if s.Type != nil {
		if len(*s.Type) > 1 {
			return nil, fmt.Errorf("list of types not supported with array constraints %#v", s)
		}
		if s.GetType()[0] != schema.TypeArray {
			return nil, fmt.Errorf("%v not supported with array constraints", s.GetType()[0])
		}
	}
	if s.UniqueItems {
		return nil, fmt.Errorf("uniqueItems are not supported")
	}
	if s.MaxItems != nil {
		return nil, fmt.Errorf("maxItems are not supported")
	}
	if s.MinItems > 0 {
		return nil, fmt.Errorf("minItems are not supported")
	}
	additionalItems := true
	if s.AdditionalItems != nil {
		if s.Items == nil {
			//any
		}
		if s.AdditionalItems.Bool != nil {
			additionalItems = *s.AdditionalItems.Bool
		}
		if !additionalItems && (s.MaxLength != nil || s.MinLength > 0) {
			return nil, fmt.Errorf("additionalItems: false and (maxItems|minItems) are not supported together")
		}
		return nil, fmt.Errorf("additionalItems are not supported")
	}
	if s.Items != nil {
		if s.Items.Object != nil {
			if s.Items.Object.Type == nil {
				//any
			} else {
				typ := s.Items.Object.GetType()[0]
				_ = typ
			}
			//TODO this specifies the type of every item in the list
		} else if s.Items.Array != nil {
			if !additionalItems {
				//TODO this specifies the length of the list as well as each ordered element's type
				//  if no type is set then any type is accepted
				maxLength := len(s.Items.Array)
				_ = maxLength
			} else {
				//TODO this specifies the types of the first few ordered items in the list
				//  if no type is set then any type is accepted
			}

		}
		return nil, fmt.Errorf("items are not supported")
	}
	return nil, nil
}
