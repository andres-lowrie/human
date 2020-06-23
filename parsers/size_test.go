package parsers

import (
	"testing"
)

func TestSizeCanParseIntoHuman(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		// At least a thousand
		{"100", false},
		// Only numbers
		{"afsdafa", false},
		// Positive
		{"1000", true},
		{"2048", true},
		{"20484046", true},
	}

	sizeP := NewSize()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := sizeP.CanParseIntoHuman(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
		})
	}
}

func TestSizeDoIntoHuman(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		// Handle "common" sizes (whatever that means lol)
		{"1", "1B"},
		{"10", "10B"},
		{"100", "100B"},
		{"1000", "1Kb"},
		{"10000", "10Kb"},
		{"100000", "100Kb"},
		{"1000000", "1Mb"},
		{"10000000", "10Mb"},
		{"100000000", "100Mb"},
		{"1000000000", "1Gb"},
		{"10000000000", "10Gb"},
		{"100000000000", "100Gb"},
		{"1000000000000", "1Tb"},
		{"10000000000000", "10Tb"},
		{"100000000000000", "100Tb"},
		{"1000000000000000", "1Pb"},
		{"10000000000000000", "10Pb"},
		{"100000000000000000", "100Pb"},
		{"1000000000000000000", "1Eb"},
		{"10000000000000000000", "10Eb"},
		{"100000000000000000000", "100Eb"},
		{"1000000000000000000000", "1Zb"},
		{"10000000000000000000000", "10Zb"},
		{"100000000000000000000000", "100Zb"},
		{"1000000000000000000000000", "1Yb"},
		{"10000000000000000000000000", "10Yb"},
		{"100000000000000000000000000", "100Yb"},
		// Should handle rouding
		{"142089140826193550568923157", "1.4Yb"},
	}

	sizeP := NewSize()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := sizeP.DoIntoHuman(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
		})
	}
}
