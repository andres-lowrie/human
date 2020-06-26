package parsers

import (
	"testing"
)

func TestNumberWordCanParseIntoHuman(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		// Only numbergroups
		{"aba", false},
		{"12af", false},
		// Lower bounds
		{"999", false},
    // Delimited numbers
    {"1,000,000", true},
    {"1.000.000", true},
    {"1 000 000", true},
    {"1a000a000", false},
		// Anything else should be parsable
		{"1000", true},
		{"100000000000", true},
		{"1338054622987", true},
	}

  numword := NewNumberWord()
  for i, tt := range tests {
    t.Run(tt.in, func(t *testing.T) {
      got := numword.CanParseIntoHuman(tt.in)
      if got != tt.out {
        t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
      }
    })
  }
}

