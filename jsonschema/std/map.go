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

package std

func Map[A any, B any](as []A, f func(A) B) []B {
	bs := make([]B, len(as))
	for i := 0; i < len(as); i++ {
		bs[i] = f(as[i])
	}
	return bs
}

func MapErr[A any, B any](as []A, f func(A) (B, error)) ([]B, error) {
	bs := make([]B, len(as))
	var err error
	for i := 0; i < len(as); i++ {
		bs[i], err = f(as[i])
		if err != nil {
			return nil, err
		}
	}
	return bs, nil
}
