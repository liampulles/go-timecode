package timecode_test

import (
	"fmt"
	"testing"

	"github.com/liampulles/go-timecode"
	"github.com/stretchr/testify/assert"
)

func TestTimecode_FormatDot(t *testing.T) {
	// Setup expectations
	var tests = []struct {
		timecode timecode.Timecode
		expected string
	}{
		{
			timecode.Zero,
			"00:00:00.000",
		},
		{
			timecode.Second,
			"00:00:01.000",
		},
		{
			timecode.Minute,
			"00:01:00.000",
		},
		{
			timecode.Hour,
			"01:00:00.000",
		},
		{
			timecode.Hour + (2 * timecode.Minute) + (3 * timecode.Second) + (4 * timecode.Millisecond),
			"01:02:03.004",
		},
		{
			-timecode.Hour + (-2 * timecode.Minute) + (-3 * timecode.Second) + (-4 * timecode.Millisecond),
			"-01:02:03.004",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			// Exercise SUT
			actual := test.timecode.FormatDot()

			// Verify result
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestTimecode_Format_NoMillisecondCase(t *testing.T) {
	// Setup fixture
	sut := timecode.Hour + (2 * timecode.Minute) + (3 * timecode.Second) + (4 * timecode.Millisecond)

	// Exercise SUT
	actual := sut.Format(false, "irrelevant")

	// Verify result
	assert.Equal(t, "01:02:03", actual)
}

func TestTimecode_FormatComma(t *testing.T) {
	// Setup fixture
	sut := timecode.Hour + (2 * timecode.Minute) + (3 * timecode.Second) + (4 * timecode.Millisecond)

	// Exercise SUT
	actual := sut.FormatComma()

	// Verify result
	assert.Equal(t, "01:02:03,004", actual)
}

func TestTimecode_String(t *testing.T) {
	// Setup fixture
	sut := timecode.Hour + (2 * timecode.Minute) + (3 * timecode.Second) + (4 * timecode.Millisecond)

	// Exercise SUT
	actual := sut.String()

	// Verify result
	assert.Equal(t, "01:02:03.004", actual)
}

func TestTimecode_WithHours(t *testing.T) {
	// Setup fixture
	sut := timecode.Hour + (2 * timecode.Minute) + (3 * timecode.Second) + (4 * timecode.Millisecond)

	// Exercise SUT
	actual := sut.WithHours(9)

	// Verify result
	assert.Equal(t, "09:02:03.004", actual.FormatDot())
}

func TestTimecode_WithMinutes(t *testing.T) {
	// Setup fixture
	sut := timecode.Hour + (2 * timecode.Minute) + (3 * timecode.Second) + (4 * timecode.Millisecond)

	// Exercise SUT
	actual := sut.WithMinutes(9)

	// Verify result
	assert.Equal(t, "01:09:03.004", actual.FormatDot())
}

func TestTimecode_WithSeconds(t *testing.T) {
	// Setup fixture
	sut := timecode.Hour + (2 * timecode.Minute) + (3 * timecode.Second) + (4 * timecode.Millisecond)

	// Exercise SUT
	actual := sut.WithSeconds(9)

	// Verify result
	assert.Equal(t, "01:02:09.004", actual.FormatDot())
}

func TestTimecode_WithMilli(t *testing.T) {
	// Setup fixture
	sut := timecode.Hour + (2 * timecode.Minute) + (3 * timecode.Second) + (4 * timecode.Millisecond)

	// Exercise SUT
	actual := sut.WithMilli(9)

	// Verify result
	assert.Equal(t, "01:02:03.009", actual.FormatDot())
}

func TestFromParams(t *testing.T) {
	// Setup expectations
	var tests = []struct {
		negative    bool
		hour        uint64
		minute      uint64
		second      uint64
		Millisecond uint64
		expected    timecode.Timecode
	}{
		{
			false,
			0, 0, 0, 0,
			timecode.Zero,
		},
		{
			false,
			0, 0, 0, 1,
			timecode.Millisecond,
		},
		{
			false,
			0, 0, 1, 0,
			timecode.Second,
		},
		{
			false,
			0, 1, 0, 0,
			timecode.Minute,
		},
		{
			false,
			1, 0, 0, 0,
			timecode.Hour,
		},
		{
			false,
			1, 2, 3, 456,
			timecode.Timecode(3723456),
		},
		{
			true,
			1, 2, 3, 456,
			timecode.Timecode(-3723456),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			// Exercise SUT
			actual := timecode.FromParams(test.negative, test.hour, test.minute, test.second, test.Millisecond)

			// Verify result
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestParse_ValidCases(t *testing.T) {
	// Setup expectations
	var tests = []struct {
		str      string
		expected timecode.Timecode
	}{
		{
			"00:00:00.000",
			timecode.Zero,
		},
		{
			"00:00:00.001",
			timecode.Millisecond,
		},
		{
			"00:00:00,001",
			timecode.Millisecond,
		},
		{
			"00:00:01.000",
			timecode.Second,
		},
		{
			"00:01:00.000",
			timecode.Minute,
		},
		{
			"01:00:00.000",
			timecode.Hour,
		},
		{
			"01:02:03.456",
			timecode.Timecode(3723456),
		},
		{
			"-01:02:03.456",
			timecode.Timecode(-3723456),
		},
		{
			"00:00:01",
			timecode.Second,
		},
		{
			"00:00:01.zzz",
			timecode.Second,
		},
		{
			"crouching.tiger.00:00:01.hidden.timecode",
			timecode.Second,
		},
		{
			"crouching.tiger.00:00:00.001.hidden.timecode",
			timecode.Millisecond,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			// Exercise SUT
			actual, err := timecode.Parse(test.str)

			// Verify result
			assert.NoError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestParse_InvalidCases(t *testing.T) {
	// Setup expectations
	var tests = []string{
		"not.a.timecode",
		"00:00:0d",
		"00:0d:00",
		"0d:00:00",
		"00:00",
		"00:00:-00",
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			// Exercise SUT
			actual, err := timecode.Parse(test)

			// Verify result
			assert.Error(t, err)
			assert.Equal(t, timecode.Zero, actual)
		})
	}
}
