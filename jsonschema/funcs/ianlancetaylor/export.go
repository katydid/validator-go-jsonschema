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

package format

// json schema validator for format: uri
func ValidateURI(s string) error {
	return uriFormat(s)
}

// json schema validator for format: iri
func ValidateIRI(s string) error {
	return iriFormat(s)
}

// json schema validator for format: uri-reference
func ValidateURIReference(s string) error {
	return uriReferenceFormat(s)
}

// json schema validator for format: iri-reference
func ValidateIRIReference(s string) error {
	return iriReferenceFormat(s)
}

// json schema validator for format: uri-template
func ValidateURITemplate(s string) error {
	return uriTemplateFormat(s)
}
