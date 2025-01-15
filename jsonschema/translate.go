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
	"fmt"
	"sort"

	"github.com/katydid/validator-go-jsonschema/validator/ast"
	"github.com/katydid/validator-go-jsonschema/validator/combinator"
)

func TranslateDraft4(schema *Schema) (*ast.Grammar, error) {
	p, err := translate(schema)
	if err != nil {
		return nil, err
	}
	return ast.NewGrammar(ast.RefLookup(map[string]*ast.Pattern{"main": p})), nil
}

func translate(schema *Schema) (*ast.Pattern, error) {
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

func translateOne(schema *Schema) (*ast.Pattern, error) {
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
		p, err := translateString(schema.String)
		return p, err
	}
	if schema.HasArrayConstraints() {
		return nil, fmt.Errorf("array not supported")
	}
	if schema.HasObjectConstraints() {
		p, err := translateObject(schema)
		return p, err
	}
	if schema.HasInstanceConstraints() {
		p, err := translateInstance(schema)
		return p, err
	}

	if len(schema.Ref) > 0 {
		return nil, fmt.Errorf("ref not supported")
	}
	if len(schema.Format) > 0 {
		return nil, fmt.Errorf("format not supported")
	}
	return ast.NewZAny(), nil
}

func translates(schemas []*Schema) ([]*ast.Pattern, error) {
	ps := make([]*ast.Pattern, len(schemas))
	for i := range schemas {
		var err error
		ps[i], err = translate(schemas[i])
		if err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func rest(xs []*ast.Pattern, index int) []*ast.Pattern {
	ys := make([]*ast.Pattern, index)
	copy(ys, xs)
	return append(ys, xs[index+1:]...)
}

func translateInstance(schema *Schema) (*ast.Pattern, error) {
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
		ps, err := translates(schema.OneOf)
		if err != nil {
			return nil, err
		}
		if len(ps) == 0 {
			return nil, fmt.Errorf("oneof of zero schemas not supported")
		}
		if len(ps) == 1 {
			return ps[0], nil
		}
		orps := make([]*ast.Pattern, len(ps))
		for i, _ := range ps {
			other := rest(ps, i)
			orps[i] = ast.NewAnd(
				ps[i],
				ast.NewNot(
					ast.NewOr(other...),
				),
			)
		}
		return ast.NewOr(orps...), nil
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

func translateType(typ SimpleType) (*ast.Pattern, error) {
	switch typ {
	case TypeArray, TypeObject:
		//This does not distinguish between arrays and objects
		return combinator.Many(combinator.InAny(combinator.Any())), nil
	case TypeBoolean:
		return combinator.Value(ast.NewType(combinator.BoolVar())), nil
	case TypeInteger:
		return combinator.Value(ast.NewType(ast.NewFunction("integer"))), nil
	case TypeNull:
		//TODO null is not being returned by json parser, but is also not empty
		return combinator.Value(combinator.Not(
			combinator.Or(
				ast.NewType(ast.NewFunction("number")),
				combinator.Or(
					ast.NewType(combinator.BoolVar()),
					ast.NewType(combinator.StringVar()),
				),
			),
		)), nil
	case TypeNumber:
		return combinator.Value(ast.NewType(ast.NewFunction("number"))), nil
	case TypeString:
		return combinator.Value(ast.NewType(combinator.StringVar())), nil
	}
	panic(fmt.Sprintf("unknown simpletype: %s", typ))
}

func translateObject(schema *Schema) (*ast.Pattern, error) {
	if schema.MaxProperties != nil {
		return nil, fmt.Errorf("maxProperties not supported")
	}
	if schema.MinProperties > 0 {
		return nil, fmt.Errorf("minProperties not supported")
	}
	required := make(map[string]struct{})
	for _, req := range schema.Required {
		required[req] = struct{}{}
	}
	requiredIf := make(map[string][]string)
	moreProperties := make(map[string]*Schema)
	if schema.Dependencies != nil {
		deps := *schema.Dependencies
		for name, dep := range deps {
			if len(dep.RequiredProperty) > 0 {
				requiredIf[name] = deps[name].RequiredProperty
			} else {
				moreProperties[name] = deps[name].Schema
			}
		}
	}
	names := []string{}
	for name, _ := range schema.Properties {
		names = append(names, name)
	}
	sort.Strings(names)
	additional := ast.NewZAny()
	if len(names) > 0 {
		nameExprs := make([]*ast.NameExpr, len(names))
		for i, name := range names {
			nameExprs[i] = ast.NewStringName(name)
		}
		additional = ast.NewZeroOrMore(
			ast.NewTreeNode(ast.NewAnyNameExcept(
				ast.NewNameChoice(nameExprs...),
			), ast.NewZAny()),
		)
	}
	if schema.AdditionalProperties != nil {
		if schema.AdditionalProperties.Bool != nil && !(*schema.AdditionalProperties.Bool) {
			additional = ast.NewEmpty()
		} else if schema.AdditionalProperties.Type != TypeUnknown {
			typ, err := translateType(schema.AdditionalProperties.Type)
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
		child, err := translate(schema.Properties[name])
		if err != nil {
			return nil, err
		}
		patterns[name] = ast.NewTreeNode(ast.NewStringName(name), child)
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
	if len(schema.PatternProperties) > 0 {
		return nil, fmt.Errorf("patternProperties not supported")
	}
	patternList := make([]*ast.Pattern, 0, len(patterns))
	for _, name := range names {

		patternList = append(patternList, patterns[name])
	}
	patternList = append(patternList, additional)
	return ast.NewInterleave(patternList...), nil
}

func optional(p *ast.Pattern) *ast.Pattern {
	return ast.NewOr(ast.NewEmpty(), p)
}

func translateNumeric(schema Numeric) (*ast.Pattern, error) {
	v := ast.NewFunction("number")
	list := []*ast.Expr{}
	notNum := combinator.Not(ast.NewType(ast.NewFunction("number")))
	if schema.MultipleOf != nil {
		mult := ast.NewFunction("multipleOf", v, combinator.DoubleConst(*schema.MultipleOf))
		list = append(list, combinator.Or(mult, notNum))
	}
	if schema.Maximum != nil {
		lt := combinator.LE(v, combinator.DoubleConst(*schema.Maximum))
		if schema.ExclusiveMaximum {
			lt = combinator.LT(v, combinator.DoubleConst(*schema.Maximum))
		}
		list = append(list, combinator.Or(lt, notNum))
	}
	if schema.Minimum != nil {
		lt := combinator.GE(v, combinator.DoubleConst(*schema.Minimum))
		if schema.ExclusiveMinimum {
			lt = combinator.GT(v, combinator.DoubleConst(*schema.Minimum))
		}
		list = append(list, combinator.Or(lt, notNum))
	}
	if len(list) == 0 {
		return combinator.Value(ast.NewType(v)), nil
	}
	return combinator.Value(and(list)), nil
}

func and(list []*ast.Expr) *ast.Expr {
	if len(list) == 0 {
		panic("unreachable")
	}
	if len(list) == 1 {
		return list[0]
	}
	return combinator.And(list[0], and(list[1:]))
}

func translateString(schema String) (*ast.Pattern, error) {
	v := combinator.StringVar()
	list := []*ast.Expr{}
	notStr := combinator.Not(ast.NewType(combinator.StringVar()))
	if schema.MaxLength != nil {
		ml := ast.NewFunction("maxLength", v, combinator.IntConst(int64(*schema.MaxLength)))
		list = append(list, combinator.Or(ml, notStr))
	}
	if schema.MinLength > 0 {
		ml := ast.NewFunction("minLength", v, combinator.IntConst(int64(schema.MinLength)))
		list = append(list, combinator.Or(ml, notStr))
	}
	if schema.Pattern != nil {
		p := combinator.Regex(combinator.StringConst(*schema.Pattern), v)
		list = append(list, combinator.Or(p, notStr))
	}
	if len(list) == 0 {
		return combinator.Value(ast.NewType(v)), nil
	}
	return combinator.Value(and(list)), nil
}

func translateArray(schema *Schema) (*ast.Pattern, error) {
	if schema.Type != nil {
		if len(*schema.Type) > 1 {
			return nil, fmt.Errorf("list of types not supported with array constraints %#v", schema)
		}
		if schema.GetType()[0] != TypeArray {
			return nil, fmt.Errorf("%v not supported with array constraints", schema.GetType()[0])
		}
	}
	if schema.UniqueItems {
		return nil, fmt.Errorf("uniqueItems are not supported")
	}
	if schema.MaxItems != nil {
		return nil, fmt.Errorf("maxItems are not supported")
	}
	if schema.MinItems > 0 {
		return nil, fmt.Errorf("minItems are not supported")
	}
	additionalItems := true
	if schema.AdditionalItems != nil {
		if schema.Items == nil {
			//any
		}
		if schema.AdditionalItems.Bool != nil {
			additionalItems = *schema.AdditionalItems.Bool
		}
		if !additionalItems && (schema.MaxLength != nil || schema.MinLength > 0) {
			return nil, fmt.Errorf("additionalItems: false and (maxItems|minItems) are not supported together")
		}
		return nil, fmt.Errorf("additionalItems are not supported")
	}
	if schema.Items != nil {
		if schema.Items.Object != nil {
			if schema.Items.Object.Type == nil {
				//any
			} else {
				typ := schema.Items.Object.GetType()[0]
				_ = typ
			}
			//TODO this specifies the type of every item in the list
		} else if schema.Items.Array != nil {
			if !additionalItems {
				//TODO this specifies the length of the list as well as each ordered element's type
				//  if no type is set then any type is accepted
				maxLength := len(schema.Items.Array)
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
