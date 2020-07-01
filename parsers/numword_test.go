package parsers

import (
	"testing"
)

func TestNumberWordCanParseIntoHuman(t *testing.T) {
	tests := []struct {
		in  string
		out bool
		err error
	}{
		// Only numbers
		{"aba", false, ErrNotANumber},
		{"12af", false, ErrNotANumber},
		{"123,afb,$$@", false, ErrNotANumber},
		// Lower/Upper Bounds
		{"999", false, ErrTooSmall},
		{"1000", true, nil},
		{"100,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000", true, nil},
		{"1,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000", false, ErrTooLarge},
		// Delimiters
		{"1,000,000", true, nil},
		{"1.000.000", true, nil},
		{"1 000 000", true, nil},
		{"1a000a000", false, ErrNotANumber},
		// Anything else should be parsable
		{"1000", true, nil},
		{"100000000000", true, nil},
		{"1338054622987", true, nil},
	}

	numword := NewNumberWord()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := numword.CanParseIntoHuman(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Error Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestNumberWordCanParseFromHumans(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		// Must be <digits> <word>
		{"1", false},
		{"million", false},
		{"one million", false},
		{"1 million!", false},
		// <word> must be in the trans table
		{"1 foo", false},
		// none of this garbage
		{"100,000 million", false},
		{"100.000 million", false},
		{"100 000 million", false},
		// These names are excluded by design
		{"1 centillion", false},
		{"1 googol", false},
		{"1 googolplex", false},

		// Tenths, Ones, Tens, Hundreds
		{"1 million", true},
		{"10 million", true},
		{"100 million", true},
		{"1.3 million", true},
		// case insensitive
		{"1 MiLlIon", true},
	}

	numword := NewNumberWord()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := numword.CanParseFromHuman(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
		})
	}
}

func TestNumberWordDoIntoHuman(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		// Each position
		{"1000", "1 thousand"},
		{"10000", "10 thousand"},
		{"100000", "100 thousand"},
		// First decimal
		{"12345678", "12.3 million"},
		// Delimiters
		{"1,000,000", "1 million"},
		{"1.000.000", "1 million"},
		{"1 000 000", "1 million"},
		// All the names
		{"1000000", "1 million"},
		{"1000000000", "1 billion"},
		{"1000000000000", "1 trillion"},
		{"1,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000", "1 vigintillion"},
	}

	numword := NewNumberWord()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := numword.DoIntoHuman(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
		})
	}
}

func TestNumberWordDoFromHuman(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"1 million", "1000000"},
		{"1.3 million", "1300000"},
		{"10 billion", "10000000000"},
		{"1 MILLION", "1000000"},
		{"100 trillion", "100000000000000"},
		{"10.3 million", "10300000"},
		{"100.3 million", "100300000"},
		{"1 vigintillion", "1000000000000000000000000000000000000000000000000000000000000000"},
	}

	numword := NewNumberWord()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := numword.DoFromHuman(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
		})
	}
}
