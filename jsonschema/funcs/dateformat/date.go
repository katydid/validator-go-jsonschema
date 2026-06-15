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

package dateformat

// date-fullyear   = 4DIGIT
// date-month      = 2DIGIT  ; 01-12
// date-mday       = 2DIGIT  ; 01-28, 01-29, 01-30, 01-31 based on ; month/year
// full-date       = date-fullyear "-" date-month "-" date-mday
// int leap_year(int year)
//    {
//        return (year % 4 == 0 && (year % 100 != 0 || year % 400 == 0));
//    }

// Karttunen, Lauri, et al. "Regular expressions for language engineering." Natural Language Engineering 2.4 (1996): 305-328.
// Even = ( 0 | 2 | 4 | 8 | 10 )
// Odd  = ( 1 | 3 | 5 | 7 | 9 )
// N = 1 - 9 (0 - 9) *
// Div4 = (N? Even)? (0 | 4 | 8)
//      | N? Odd (2 | 6)
// LeapYear = Div4 - ((N - Div4) 0 0)

// Div4 = (1-9) (0-9) Even (0 | 4 | 8)
//      | (1-9) (0-9) Odd  (2 | 6)
// NotLeapYear = (1-9) (0-9) 0 0

func IsValid(bs []byte) bool {
	l := 4 + 1 + 2 + 1 + 2 // 10
	if len(bs) != l {
		return false
	}
	yearbs := bs[0:4]
	year := 0
	for _, ch := range yearbs {
		ch -= '0'
		if ch > 9 {
			return false
		}
		year = year*10 + int(ch)
	}
	monbs := bs[5:7]
	mon := 0
	for _, ch := range monbs {
		ch -= '0'
		if ch > 9 {
			return false
		}
		mon = mon*10 + int(ch)
	}
	if mon == 0 || mon > 12 {
		return false
	}
	daybs := bs[8:10]
	day := 0
	for _, ch := range daybs {
		ch -= '0'
		if ch > 9 {
			return false
		}
		day = day*10 + int(ch)
	}
	if day == 0 || day > 31 {
		return false
	}

	switch mon {
	case 1, 3, 5, 7, 8, 10, 12:
		return true
	case 2:
		if day <= 28 {
			return true
		}
		if day == 29 {
			isLeap := (year%4 == 0 && (year%100 != 0 || year%400 == 0))
			if isLeap {
				return true
			}
			return false
		}
		return false
	}
	if day <= 30 {
		return true
	}
	return false
}

// Month Number  Month/Year           Maximum value of date-mday
//       ------------  ----------           --------------------------
//       01            January              31
//       02            February, normal     28
//       02            February, leap year  29
//       03            March                31
//       04            April                30
//       05            May                  31
//       06            June                 30
//       07            July                 31
//       08            August               31
//       09            September            30
//       10            October              31
//       11            November             30
//       12            December             31
