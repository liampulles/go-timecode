package timecode

import (
	"fmt"
	"regexp"
	"strconv"
)

// Regex can be used to validate timecodes and capture the sign (optional)
// hours, minutes, seconds, and milliseconds (optional) groups. Digits are
// properly checked to fall under a 24 hour clock, and the millisecond seperator
// may either be a dot or a comma.
var Regex = regexp.MustCompile(`([-])?([01]\d|2[0123]):([012345]\d):([012345]\d)(?:[.,](\d{3}))?`)

// Some basic Timecode units
const (
	Zero        = Timecode(0)
	Millisecond = Timecode(1)
	Second      = Millisecond * 1000
	Minute      = Second * 60
	Hour        = Minute * 60
)

// Timecode defines a timecode that would typically be used to denote a duration
// or position in media, e.g. for a subtitle or audio track.
//
// Timecode has millisecond resolution, and may be negative.
//
// Timecodes may be directly added/subtracted, and multiplied by factors (as
// long as the result is cast back to Timecode)
type Timecode int64

// Check we implement the interface
var _ fmt.Stringer = Timecode(0)

// HourMinuteSecondMilli returns the constituent elements of a timecode
func (t Timecode) HourMinuteSecondMilli() (int64, int64, int64, int64) {
	i := int64(t)

	milli := i % 1000
	i = i / 1000
	second := i % 60
	i = i / 60
	minute := i % 60
	i = i / 60
	hour := i

	return hour, minute, second, milli
}

// Format formats a Timecode into a string. If withMilli is true,
// then milliSeperator is used to separate the seconds section from the
// milliseconds section.
func (t Timecode) Format(withMilli bool, milliSeperator string) string {
	tWip := t

	negative := false
	if tWip < 0 {
		negative = true
		tWip = tWip * -1
	}

	hour, minute, second, milli := tWip.HourMinuteSecondMilli()

	var result string
	if withMilli {
		result = fmt.Sprintf("%02d:%02d:%02d%s%03d",
			hour, minute, second, milliSeperator, milli)
	} else {
		result = fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
	}
	if negative {
		return fmt.Sprintf("-%s", result)
	}
	return result
}

// FormatDot will format as e.g. 01:02:03.004
func (t Timecode) FormatDot() string {
	return t.Format(true, ".")
}

// FormatComma will format as e.g. 01:02:03,004
func (t Timecode) FormatComma() string {
	return t.Format(true, ",")
}

// String is the same as FormatDot. It should be used for logging; non-business
// logic purposes as this format is NOT guaranteed.
func (t Timecode) String() string {
	return t.FormatDot()
}

// FromParams constructs a Timecode from its constituent parts. Values may be
// negative or positive, but they should all have the same sign to avoid
// unexpected results. E.g. to create a timecode like -01:02:03.456, use
// FromParams(-1,-2,-3,-456).
func FromParams(hour, minute, second, milli int64) Timecode {
	total := milli
	total += second * 1000
	total += minute * 1000 * 60
	total += hour * 1000 * 60 * 60
	return Timecode(total)
}

// Parse extracts a Timecode from a string. The following are valid examples
// of Timecodes (non-valid Timecodes will return an error):
// "01:02:03.456"
// "01:02:03,456"
// "-01:02:03.456"
// "01:02:03"
// "crouching.tiger.01:02:03.456.hidden.timecode"
func Parse(str string) (Timecode, error) {
	m := Regex.FindStringSubmatch(str)
	if len(m) == 0 {
		return Timecode(0), fmt.Errorf("not a timecode")
	}

	mult := int64(1)
	if isNotEmpty(m, 1) {
		mult = -1
	}

	hour := parseNumber(m, 2)
	minute := parseNumber(m, 3)
	second := parseNumber(m, 4)
	milli := parseNumber(m, 5)

	return FromParams(mult*hour, mult*minute, mult*second, mult*milli), nil
}

func parseNumber(regexMatch []string, i int) int64 {
	value := regexMatch[i]
	result, _ := strconv.ParseInt(value, 10, 64)
	return result
}

func isNotEmpty(regexMatch []string, i int) bool {
	value := regexMatch[i]
	return value != ""
}
