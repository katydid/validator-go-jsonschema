# uniqueItems is unsupported

This document explains why we chose not to support the uniqueItems operator from JSON Schema.

tl;dr It is not possible to add the `uniqueItems` operator inside Katydid's optimized derivative algorithm.

*The document assumes a familiarity with Brzozowski's derivatives for regular expressions.* 

## What is uniqueItems

The `uniqueItems` operator makes it possible to specify that an `array` is really a Set.

Given the following JSON Schema:

```json
{"uniqueItems": true}
```

We expect it to match arrays that only have uniqueItems, for example:

```json
[]
[1, 2]
[1, 2, 3]
```

We expect it to **not** match arrays with duplicate items, for example:

```json
[1, 1]
[1, 2, 3, 2]
```

This is not only true for simple values, but also include objects and arrays.
For example, we expect the following to also **not** match:

```json
[["foo"], ["bar"], ["foo"]]
[{"a": 1, "b": 2}, {"b": 2, "a": 1}]
```

`uniqueItems` is an operator that can be composed with other json schemas, for example:

```json
{
  "type": "array",
  "items": {
    "type": "string",
    "pattern": "^\\d{3}-\\d{3}-\\d{4}$"
  },
  "uniqueItems": false
}
```

## Theoretical Considerations

All JSONSchemas, without the `uniqueItems` operator, can be expressed using Monadic Second Order Logic (MSO) and Tree Automata.
Second-order operators are required to handle definitions, references and regular expressions, but all other operators are expressible using First-Order Logic.
The `{"uniqueItems":true}` constraint is the only exception, which is not expressible using MSO.

For example, with the `uniqueItems` operator it is possible to specify perfectly balanced binary trees, which are not expressible with Tree Automata:

```json
{ "definitions": 
  { "Binary":
    { "anyOf":
      [
        { 
          "enum": ["leaf"]
        },
        { 
          "type": "array",
          "not": { 
            "uniqueItems": true 
          },
          "minItems": 2,
          "maxItems": 2,
          "items": { 
            "$ref": "#/definitions/Binary" 
          }
        }
      ]
    } 
  }
}
```

The only JSON values that will match this JSONSchema is either a `"leaf"` string or an array of two items that are duplicates of each other:

```json
"leaf"
["leaf", "leaf"]
[["leaf", "leaf"], ["leaf", "leaf"]]
...
```

## Performance Considerations

Validating `uniqueItems` also adds a performance overhead as we either need to:
  * sort the array before validation, which no longer makes the processing time linear and makes streaming processing unviable, or
  * use a hashtable to keep track of the items, which no longer makes the space required constant, but results in linear processing time.

Currently the algorithm processes in linear time and after enough memoization, there is no more memory allocation required.
The `uniqueItems` operator requires us to lose our gauranteed linear processing time or allocate memory, which will either way make the algorithm slower.

Another performance overhead is the cache hits would be consideribly lower, as each expression is now less likely to be a duplicate of another expression.
This happens because the derivative uses the input character in the result, which means that each `uniqueItems` derivatives now includes random characters from the input.

## Naive Derivative Algorithm

Katydid's algorithm is based on Brzozowski's derivatives.
The only other derivative based JSON Schema validator (at the time of writing) created by Mary Holstege in 2019, supports uniqueItems, but does so outside of the derivative based algorithm.
Even though `uniqueItems` is a bad fit for the derivative algorithm, it does not mean that it is not possible.

Let us consider the string case:

```lean
nullable : Regex -> Bool
...
nullable UniqueItems = true

derive : Char -> Regex -> Regex
derive ...
derive c (Char a) = if c == a then EmptyStr else EmptySet
derive ...
derive c UniqueItems = (And UniqueItems (Not (Concat (Not EmptySet) (Concat (Char c) (Not EmptySet))))) -- (uniqueItems&!(.*c.*))
```

The derivative of `uniqueItems` is `uniqueItems` and a constraint that the character does not appear in the string again: `!(.*c.*)`.
There are other ways to accomplish the same thing, for example, uniqueItems could store a list of already seen items:

```lean
derive c (UniqueItems []) = UniqueItems [c]
derive c (UniqueItems xs) = if contains xs c then EmptySet else UniqueItems (c::xs)
```

Unfortunately this does not influence the reasoning in the rest of this document.

This is much harder in the case of trees (arrays and objects in JSON), but would be possible in a naive derivative algorithm.
Unfortunately Katydid is not a naive derivative algorithm and has some extra contraints and assumptions, which makes this approach very hard.

## Optimized Derivative Algorithm

Derivatives for regular expressions are memoizable into Deterministic Finite Automata (DFA).
If we assume a regex is a state, then we can transition based on a character to another regex (or state):

```lean
derive : Char -> Regex -> Regex
```

Problem is a unicode `Char` is quite a large alphabet to memoize on.
We could pass in every unicode character to every regex and keep track of visited ones and so explore the whole state space, but this will create a ginormous table.
When we use trees (objects, arrays) as our input alphabet, as Katydid does, then this alphabet becomes infinite, which means this technique becomes impossible to use.
In practice, different techniques are used.

One such option is to abstract away from `Char` operators and rather use predicates.
This is called Symbolic Regular Expressions.
This means instead of having `a|b`, we have an expression:

```
(\x => x == a)|(\x => x == b)
```

Katydid is based on derivatives for Symbolic Regular Expressions and allows users to even create their own predicates.
Given these predicates, we know that the input character either matches the predicate or it does not.
This means we can explore the possible state space without any inputs.
We do this by calculating a derivative "if expression" for a regular expression.

```lean
deriveIfExpr : Regex -> IfExpr
...
deriveIfExpr (Pred pred) = IfExpr pred EmptyStr EmptySet
deriveIfExpr (Or x y) = 
    match (deriveIfExpr x, deriveIfExpr y) with
    | (IfExpr cndx thnx elsx, IfExpr cndy thny elsy) =>
        IfExpr cndx (IfExpr cndy (Or thnx thny) (Or thnx elsy)) (IfExpr cndy (Or elsx thny) (Or elsx elsy))
...
```

If we repeat this process and keep track of visited states, we can explore the whole state space.
Notice we do not need any input to calculate all possibile resulting regular expressions.

The `uniqueItems` operator is different and breaks the assumption that we can explore the whole state space without any input.
This operator requires the input to calculate the next derivative.

```
derive c (UniqueItems []) = UniqueItems [c]
derive c (UniqueItems xs) = if contains xs c then EmptySet else UniqueItems (c::xs)
```

It cannot know what the next derivative is without the input `c`.

The symbolic derivative algorithm in Katydid uses a memoization strategy that is built on the assumption that we can derive "if expressions" from our Katydid expressions.
This means that to support `uniqueItems`, we would need to remove this assumption and create a totally new algorithm.

## Conclusion

The theoretical and performance considerations are both things that can be worked around.
We can make Katydid's validator extendible with new operators and not push these issues into the core algorithm.
Then embracing these issues will be a concious decision by the extender, for example this json schema validator.
We even created a [proof of concept for extendible operators](https://github.com/katydid/validator-go/pull/37), but this unfortunately won't be enough.
The `uniqueItems` is special in another way, which breaks even more of the assumptions that the Katydid algorithm is built on.

It is possible to use `uniqueItems` in a derivative algorithm, but it is **not** possible to use `uniqueItems` in Katydid's optimized derivative algorithm.
This is because `uniqueItems` breaks the assumption that we can calculate an "if expression" for all possible states, without providing any input and without exploring the whole input alphabet.

## References

* [Brzozowski, Janusz A. "Derivatives of regular expressions." Journal of the ACM (JACM) 11.4 (1964): 481-494.](https://dl.acm.org/doi/abs/10.1145/321239.321249)
* [JSON-Schema-Test-Suite/tests/draft7/uniqueItems.json](https://github.com/json-schema-org/JSON-Schema-Test-Suite/blob/83e866b46c9f9e7082fd51e83a61c5f2145a1ab7/tests/draft7/uniqueItems.json#L134)
* [A Tour of JSON Schema - Unique Array Items](https://tour.json-schema.org/content/04-Arrays/02-Unique-Items)
* [Application of Brzozowski Derivatives to JSON Schema Validation - Holstege, Mary - 2019](https://doi.org/10.4242/BalisageVol23.Holstege01)
* [Stanford, Caleb, Margus Veanes, and Nikolaj Bjørner. "Symbolic Boolean derivatives for efficiently solving extended regular expression constraints." Proceedings of the 42nd ACM SIGPLAN International Conference on Programming Language Design and Implementation. 2021.](https://dl.acm.org/doi/pdf/10.1145/3453483.3454066)