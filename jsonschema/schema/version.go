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

package schema

import (
	"math"
	"strings"
)

type Version int

const VersionUnknown = 0
const VersionDraft4 = 4
const VersionDraft6 = 6
const VersionDraft7 = 7
const VersionDraft2019 = 2019
const VersionDraft2020 = 2020
const VersionLatest = math.MaxInt

var strToVersion = map[string]Version{
	"https://json-schema.org/schema":               VersionLatest,
	"https://json-schema.org/draft/2020-12/schema": VersionDraft2020,
	"http://json-schema.org/draft/2019-09/schema":  VersionDraft2019,
	"http://json-schema.org/draft-07/schema":       VersionDraft7,
	"http://json-schema.org/draft-06/schema":       VersionDraft6,
	"http://json-schema.org/draft-04/schema":       VersionDraft4,
}

var versionToStr = map[Version]string{}

func init() {
	for k, v := range strToVersion {
		versionToStr[v] = k
	}
}

func detectVersion(url string) Version {
	u := strings.Split(url, "#")[0]
	strings.Replace(u, "http://", "https://", 1)
	return strToVersion[u]
}

func setDefaultVersion(s *Schema, defaultVersion Version) {
	defaultVersionStr := versionToStr[defaultVersion]
	s.Walk(func(sch *Schema) {
		if detectVersion(sch.Schema) == VersionUnknown {
			sch.Schema = defaultVersionStr
		}
	})
}
