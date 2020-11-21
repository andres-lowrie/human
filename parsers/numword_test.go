package parsers

import (
	"testing"
)

func TestNumberWordCanParseFromMachine(t *testing.T) {
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
			got, err := numword.CanParseFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Error Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestNumberWordCanParseIntoMachine(t *testing.T) {
	tests := []struct {
		in  string
		out bool
		err error
	}{
		// Must be <digits> <word>
		{"1", false, ErrNotADigitWordCombo},
		{"million", false, ErrNotADigitWordCombo},
		{"one million", false, ErrNotADigitWordCombo},
		{"1 million!", false, ErrNotADigitWordCombo},
		// <word> must be in the trans table
		{"1 foo", false, ErrNotADigitWordCombo},
		// none of this garbage
		{"100,000 million", false, ErrNotADigitWordCombo},
		{"100.000 million", false, ErrNotADigitWordCombo},
		{"100 000 million", false, ErrNotADigitWordCombo},
		// These names are excluded by design
		{"1 centillion", false, ErrNotADigitWordCombo},
		{"1 googol", false, ErrNotADigitWordCombo},
		{"1 googolplex", false, ErrNotADigitWordCombo},

		// Tenths, Ones, Tens, Hundreds
		{"1 million", true, nil},
		{"10 million", true, nil},
		{"100 million", true, nil},
		{"1.3 million", true, nil},
		// case insensitive
		{"1 MiLlIon", true, nil},
	}

	numword := NewNumberWord()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := numword.CanParseIntoMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestNumberWordDoFromMachine(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err error
	}{
		// Each position
		{"1000", "1 thousand", nil},
		{"10000", "10 thousand", nil},
		{"100000", "100 thousand", nil},
		// First decimal
		{"12345678", "12.3 million", nil},
		// Delimiters
		{"1,000,000", "1 million", nil},
		{"1.000.000", "1 million", nil},
		{"1 000 000", "1 million", nil},
		// All the names
		{"1000000", "1 million", nil},
		{"1000000000", "1 billion", nil},
		{"1000000000000", "1 trillion", nil},
		{"1,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000,000", "1 vigintillion", nil},
	}

	numword := NewNumberWord()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := numword.DoFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestNumberWordDoIntoMachine(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err error
	}{
		{"1 million", "1000000", nil},
		{"1.3 million", "1300000", nil},
		{"10 billion", "10000000000", nil},
		{"1 MILLION", "1000000", nil},
		{"100 trillion", "100000000000000", nil},
		{"10.3 million", "10300000", nil},
		{"100.3 million", "100300000", nil},
		{"1 vigintillion", "1000000000000000000000000000000000000000000000000000000000000000", nil},
	}

	numword := NewNumberWord()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := numword.DoIntoMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}
