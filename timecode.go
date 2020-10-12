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
var _ fmt.Stringer = Zero

// HourMinuteSecondMilli returns the constituent elements of a timecode.
func (t Timecode) HourMinuteSecondMilli() (uint64, uint64, uint64, uint64) {
	i := uint64(t)
	if t.IsNegative() {
		i = uint64(t * -1)
	}

	milli := i % 1000
	i = i / 1000
	second := i % 60
	i = i / 60
	minute := i % 60
	i = i / 60
	hour := i

	return hour, minute, second, milli
}

// IsNegative is true if Timecode is below zero
func (t Timecode) IsNegative() bool {
	return t < Zero
}

// WithHours returns a new Timecode with the hours set as given.
func (t Timecode) WithHours(hour uint64) Timecode {
	_, m, s, ms := t.HourMinuteSecondMilli()
	return FromParams(t.IsNegative(), hour, m, s, ms)
}

// WithMinutes returns a new Timecode with the minutes set as given.
func (t Timecode) WithMinutes(minutes uint64) Timecode {
	h, _, s, ms := t.HourMinuteSecondMilli()
	return FromParams(t.IsNegative(), h, minutes, s, ms)
}

// WithSeconds returns a new Timecode with the seconds set as given.
func (t Timecode) WithSeconds(seconds uint64) Timecode {
	h, m, _, ms := t.HourMinuteSecondMilli()
	return FromParams(t.IsNegative(), h, m, seconds, ms)
}

// WithMilli returns a new Timecode with the milliseconds set as given.
func (t Timecode) WithMilli(milli uint64) Timecode {
	h, m, s, _ := t.HourMinuteSecondMilli()
	return FromParams(t.IsNegative(), h, m, s, milli)
}

// Format formats a Timecode into a string.
//
// If withMilli is true, then milliSeperator is used to separate the seconds
// section from the milliseconds section.
func (t Timecode) Format(withMilli bool, milliSeperator string) string {
	negative := t.IsNegative()
	h, m, s, ms := t.HourMinuteSecondMilli()

	var result string
	if withMilli {
		result = fmt.Sprintf("%02d:%02d:%02d%s%03d",
			h, m, s, milliSeperator, ms)
	} else {
		result = fmt.Sprintf("%02d:%02d:%02d", h, m, s)
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

// FromParams constructs a Timecode from its constituent parts.
func FromParams(negative bool, hour, minute, second, milli uint64) Timecode {
	total := Timecode(milli) * Millisecond
	total += Timecode(second) * Second
	total += Timecode(minute) * Minute
	total += Timecode(hour) * Hour
	if negative {
		total *= -1
	}
	return total
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
		return Zero, fmt.Errorf("[%s] is not a timecode", str)
	}

	negative := isNotEmpty(m, 1)
	hour := parseNumber(m, 2)
	minute := parseNumber(m, 3)
	second := parseNumber(m, 4)
	milli := parseNumber(m, 5)

	return FromParams(negative, hour, minute, second, milli), nil
}

func parseNumber(regexMatch []string, i int) uint64 {
	value := regexMatch[i]
	result, _ := strconv.ParseUint(value, 10, 64)
	return result
}

func isNotEmpty(regexMatch []string, i int) bool {
	value := regexMatch[i]
	return value != ""
}
