package rest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestServerModeFromEnv(t *testing.T) {
	cases := []struct {
		name     string
		env      string
		expected string
	}{
		{
			name:     "production",
			env:      "production",
			expected: ServerModeProduction,
		},
		{
			name:     "staging",
			env:      "staging",
			expected: ServerModeStaging,
		},
		{
			name:     "development",
			env:      "development",
			expected: ServerModeDevelopment,
		},
		{
			name:     "unexpected value",
			env:      "foo",
			expected: ServerModeDevelopment,
		},
		{
			name:     "empty string",
			env:      "",
			expected: ServerModeDevelopment,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := serverModeFromEnv(c.env)
			if diff := cmp.Diff(c.expected, result); diff != "" {
				t.Errorf("Mode value is not valid (-expected +actual):\n%s", diff)
			}
		})
	}
}
