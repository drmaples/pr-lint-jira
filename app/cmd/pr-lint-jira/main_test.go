package main

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultTitleBodyRegex(t *testing.T) {
	re := regexp.MustCompile(defaultTitleBodyRegex)

	tests := []struct {
		input    string
		expected bool
	}{
		{"foo", false},                     // do better
		{"FOO-123", false},                 // no brackets
		{"[f-123]", false},                 // min 2 char slug
		{"[hi-]", false},                   // need digits after slug
		{"nothing in here", false},         // no tix at all
		{"what [is-going] on here", false}, // no tix at all
		{"[123-ABC]", false},               // reversed format
		{"[-123]", false},                  // no slug
		{"[]", false},                      // empty brackets
		{"[ABCD-ABC]", false},              // letters instead of numbers
		{"{ABCD-123}", false},              // wrong bracket type

		{"[HI-7]", true},
		{"[FOO-123]", true},
		{"[zzz-123]", true},
		{"[zZz-123]", true},
		{"[AB-0]", true},
		{"[AB-1]", true},
		{"[AB-999999]", true},
		{"[XXXXXXXXXX-0000000000]", true},
		{"start junk [zZz-123]", true},
		{"[zZz-123] end junk", true},
		{"start junk [zZz-123] end junk", true},
		{"[FOO-123] [BAR-456]", true},         // multiple tix ok
		{"[ABC-1] some text [XYZ-999]", true}, // multiple tix ok
		{"[f-] some text [XYZ-999]", true},    // at least 1 good tix
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := re.MatchString(tt.input)
			assert.Equal(t, tt.expected, got, "regex.MatchString(%q)", tt.input)
		})
	}
}

func TestRuns_with_missing_required_fields(t *testing.T) {
	err := run(context.TODO())
	require.Error(t, err)
}

func TestRuns_with_context_cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	require.Error(t, run(ctx))
}
