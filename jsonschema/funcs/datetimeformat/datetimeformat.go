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

package datetimeformat

// https://datatracker.ietf.org/doc/html/rfc3339#section-5.6
// https://datatracker.ietf.org/doc/html/rfc4234

// date-fullyear   = 4DIGIT
// _digit         : '0' - '9';
// _4digit        : _digit _digit _digit _digit
//                ;
// _datefullyear  : _4digit
//                ;

// date-month      = 2DIGIT  ; 01-12
// date-mday       = 2DIGIT  ; 01-28, 01-29, 01-30, 01-31 based on
//                             ; month/year
//  Month Number  Month/Year           Maximum value of date-mday
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
// int leap_year(int year)
//    {
//        return (year % 4 == 0 && (year % 100 != 0 || year % 400 == 0));
//    }
// time-hour       = 2DIGIT  ; 00-23
// Although ISO 8601 permits the hour to be "24", this profile of ISO
//    8601 only allows values between "00" and "23"
// time-minute     = 2DIGIT  ; 00-59
// time-second     = 2DIGIT  ; 00-58, 00-59, 00-60 based on leap second rules
// The grammar element time-second may have the value "60" at the end of
//    months in which a leap second occurs -- to date: June (XXXX-06-
//    30T23:59:60Z) or December (XXXX-12-31T23:59:60Z)
// time-secfrac    = "." 1*DIGIT
// time-numoffset  = ("+" / "-") time-hour ":" time-minute
// time-offset     = "Z" / time-numoffset

// partial-time    = time-hour ":" time-minute ":" time-second
//                     [time-secfrac]
// full-date       = date-fullyear "-" date-month "-" date-mday
// full-time       = partial-time time-offset

// date-time       = full-date "T" full-time

func bytesToInt(bs []byte) (int, bool) {
	res := 0
	for _, ch := range bs {
		ch -= '0'
		if ch > 9 {
			return 0, false
		}
		res = res*10 + int(ch)
	}
	return res, true
}

func validInt(bs []byte) bool {
	for _, ch := range bs {
		ch -= '0'
		if ch > 9 {
			return false
		}
	}
	return true
}

func isValidDay(year, mon, day int) bool {
	if mon == 0 || mon > 12 {
		return false
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

func IsValid(bs []byte) bool {
	if len(bs) < 20 {
		return false
	}

	yearbs := bs[0:4]
	year, ok := bytesToInt(yearbs)
	if !ok {
		return false
	}
	if bs[4] != '-' {
		return false
	}
	monbs := bs[5:7]
	mon, ok := bytesToInt(monbs)
	if !ok {
		return false
	}
	if bs[7] != '-' {
		return false
	}
	daybs := bs[8:10]
	day, ok := bytesToInt(daybs)
	if !ok {
		return false
	}
	if !isValidDay(year, mon, day) {
		return false
	}
	if bs[10] != 'T' && bs[10] != 't' {
		return false
	}
	hour, ok := bytesToInt(bs[11:13])
	if !ok { // hour
		return false
	}
	if hour > 23 {
		return false
	}
	if bs[13] != ':' {
		return false
	}
	min, ok := bytesToInt(bs[14:16])
	if !ok { // min
		return false
	}
	if bs[16] != ':' {
		return false
	}
	if min > 59 {
		return false
	}
	sec, ok := bytesToInt(bs[17:19])
	if !ok {
		return false
	}
	if sec > 60 {
		return false
	}
	i := 19
	// optional time-secfrac
	if bs[19] == '.' {
		i++
		for bs[i] != 'Z' && bs[i] != 'z' && bs[i] != '+' && bs[i] != '-' {
			switch bs[i] {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			default:
				return false
			}
			i++
		}
		if i == 20 {
			return false
		}
	}
	if len(bs) == i+1 && (bs[i] == 'Z' || bs[i] == 'z') {
		if sec == 60 && !isLeapSecond(year, mon, day, hour, min) {
			return false
		}
		return true
	}
	if len(bs) != i+6 {
		return false
	}
	if bs[i] != '+' && bs[i] != '-' {
		return false
	}
	offsethour, ok := bytesToInt(bs[i+1 : i+3])
	if !ok { // offset hour
		return false
	}
	if offsethour > 23 {
		return false
	}
	if bs[i+3] != ':' {
		return false
	}
	offsetmin, ok := bytesToInt(bs[i+4 : i+6])
	if !ok { // offset min
		return false
	}
	if offsetmin > 59 {
		return false
	}
	if sec == 60 && !isLeapSecond(year, mon, day, hour+offsethour, min+offsetmin) {
		return false
	}
	return true
}

func isLeapSecond(year, mon, day, hour, min int) bool {
	if hour != 23 {
		return false
	}
	if min != 59 {
		return false
	}
	if mon == 12 {
		if day != 31 {
			return false
		}
		return true
	}
	if mon == 6 {
		if day != 30 {
			return false
		}
		return true
	}
	return false
}
