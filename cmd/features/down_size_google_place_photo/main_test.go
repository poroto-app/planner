package main

import (
	"github.com/google/go-cmp/cmp"
	"poroto.app/poroto/planner/internal/domain/utils"
	"testing"
)

func TestRewriteUrl(t *testing.T) {
	cases := []struct {
		name     string
		url      string
		width    int
		height   int
		expected *string
	}{
		{
			name:     "valid url",
			url:      "https://lh3.googleusercontent.com/places/photo_id=s1600-w2000-h2000",
			width:    500,
			height:   500,
			expected: utils.ToPointer("https://lh3.googleusercontent.com/places/photo_id=s1600-w500-h500"),
		},
		{
			name:     "invalid url",
			url:      "https://lh3.googleusercontent.com/places/photo_id",
			width:    500,
			height:   500,
			expected: nil,
		},
		{
			name:     "replace only last part",
			url:      "https://lh3.googleusercontent.com/places/photo_w2000-h2000=s1600-w2000-h2000",
			width:    500,
			height:   500,
			expected: utils.ToPointer("https://lh3.googleusercontent.com/places/photo_w2000-h2000=s1600-w500-h500"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, _ := rewriteUrl(c.url, c.width, c.height)
			if diff := cmp.Diff(result, c.expected); diff != "" {
				t.Errorf("expected: %v, result: %v", c.expected, result)
			}
		})
	}
}
