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

package translate

import "github.com/katydid/validator-go/validator/ast"

func translateFormat(format string) (*ast.Expr, error) {
	switch format {
	case "date":
		return dateExpr(), nil
	case "date-time":
		return datetimeExpr(), nil
	case "email":
		return emailExpr(), nil
	case "hostname":
		return hostNameExpr(), nil
	case "json-pointer":
		return jsonPointerExpr(), nil
	case "relative-json-pointer":
		return relativeJSONPointerExpr(), nil
	case "uuid":
		return uuidExpr(), nil
	case "duration":
		return durationExpr(), nil
	case "ipv4":
		return ipv4Expr(), nil
	case "ipv6":
		return ipv6Expr(), nil
	case "time":
		return timeExpr(), nil
	case "uri", "iri":
		return uriExpr(), nil
	case "uri-reference", "iri-reference":
		return uriReferenceExpr(), nil
	case "uri-template":
		return uriTemplateExpr(), nil
	case "period":
		return periodExpr(), nil
	case "semver":
		return semverExpr(), nil
	default:
		// A format attribute can generally only validate a given set of instance types.
		// If the type of the instance to validate is not in this set, validation for this format attribute and instance SHOULD succeed.
		return anyExpr(), nil
	}
}
