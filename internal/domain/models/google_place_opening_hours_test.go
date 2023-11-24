package models

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParseTimeString(t *testing.T) {
	cases := []struct {
		name     string
		timeStr  string
		expected TimeHHMM
	}{
		{
			name:    "valid time string",
			timeStr: "0000",
			expected: TimeHHMM{
				Hour:   0,
				Minute: 0,
			},
		},
		{
			name:    "valid time string",
			timeStr: "2359",
			expected: TimeHHMM{
				Hour:   23,
				Minute: 59,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			timeHHMM, err := parseTimeString(c.timeStr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(*timeHHMM, c.expected); diff != "" {
				t.Errorf("parseTimeString() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
