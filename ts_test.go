package main

import (
	"testing"
)

func TestFormatTimestamp(t *testing.T) {
	data := map[Timestamp]string{
		((2 * 60) + 23) * 1000:        "02:23.000",
		(3600 + (2 * 60) + 23) * 1000: "01:02:23.000",
		123:                           "00:00.123",
	}
	for k, v := range data {
		if out := k.String(); out != v {
			t.Fatalf("Output mismatch, should be: \"%s\", is: %s", v, out)
		}
	}
}

func TestGetMs(t *testing.T) {
	data := map[string]Timestamp{
		"10":       -1,
		"10.220":   -1,
		"2:30.123": 150123,
		"1:20":     80000,
		"1:00:00":  3600000,
	}

	for k, v := range data {
		if res, err := getMs(k); v != res || err != nil {
			if err != nil {
				if res != -1 {
					t.Fatalf("TestGetMs Error: %s", err.Error())
				}
			} else {
				t.Fatalf("Error converting string to [ms] (%v). Expected (%s)->(%v), got: %v", err, k, v, res)
			}
		}

	}
}
