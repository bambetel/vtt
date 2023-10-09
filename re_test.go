package main

import (
	"regexp"
	"testing"
)

func TestFormatTimestamp(t *testing.T) {
	data := map[int]string{
		((2 * 60) + 23) * 1000:        "02:23.000",
		(3600 + (2 * 60) + 23) * 1000: "01:02:23.000",
		123:                           "00:00.123",
	}
	for k, v := range data {
		if out := formatTimestamp(k); out != v {
			t.Fatalf("Output mismatch, should be: \"%s\", is: %s", v, out)
		}
	}
}

func TestGetMs(t *testing.T) {
	data := map[string]int{
		"10":       -1,
		"10.220":   -1,
		"2:30.123": 150123,
		"1:20":     80000,
		"1:00:00":  3600000,
	}

	for k, v := range data {
		if res := getMs(k); res != v {
			t.Fatalf("Error converting string to [ms]. Expected (%s)->(%d), got: %d", k, v, res)
		}
	}
}

func TestRegexTimestamps(t *testing.T) {
	validTSs := []string{"1:30:2", "00:0231.210"}
	invalidTSs := []string{"9.20", "10m12s", "10:20-2.01.210", "10:20m23 --> 02:10.00", "10:20.300 --> 02:10m.00"}

	reTSLine := regexp.MustCompile("^\\d+:[^\\s-]*\\d+$")
	for _, v := range validTSs {
		if !reTSLine.MatchString(v) {
			t.Fatalf("String should be matched as a timestamp: %s", v)
		}
	}

	for _, v := range invalidTSs {
		if reTSLine.MatchString(v) {
			t.Fatalf("String shouldn't be matched as a timestamp: %s", v)
		}
	}
}

func TestRegexCues(t *testing.T) {
	reCueLine := regexp.MustCompile("^(\\d+:\\d+\\S*)\\s*-+>?\\s*(\\S*:\\d+(\\.\\d+)?)$")
	// just beginning time
	// reTSLine := regexp.MustCompile("^\\d+:[^\\s-]*\\d+$")

	validTSs := []string{"1:30-10:2", "00:01-00:01.210"}
	invalidTSs := []string{"9.20", "10-120", "10:20-2.01.210", "10:20m23 --> 02:10.00", "10:20.300 --> 02:10m.00"}

	for _, v := range validTSs {
		if !reCueLine.MatchString(v) {
			t.Fatalf("String should be matched as cue: %s", v)
		}
	}

	for _, v := range invalidTSs {
		if reCueLine.MatchString(v) {
			t.Fatalf("String shouldn't be matched as cue: %s", v)
		}
	}
}
