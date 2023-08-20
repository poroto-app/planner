package plangen

import "testing"

func TestParseTimeString(t *testing.T) {
	cases := []struct {
		name           string
		timeStr        string
		expectedHour   int
		expectedMinute int
	}{
		{
			name:           "valid time string",
			timeStr:        "0000",
			expectedHour:   0,
			expectedMinute: 0,
		},
		{
			name:           "valid time string",
			timeStr:        "2359",
			expectedHour:   23,
			expectedMinute: 59,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hour, minute, err := parseTimeString(c.timeStr)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if hour != c.expectedHour {
				t.Errorf("expected: %v, actual: %v", c.expectedHour, hour)
			}
			if minute != c.expectedMinute {
				t.Errorf("expected: %v, actual: %v", c.expectedMinute, minute)
			}
		})
	}
}
