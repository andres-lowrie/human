package parsers

import (
	"testing"
)

func TestNumberGroupCanParseIntoHuman(t *testing.T) {
	tests := []struct {
		in  string
		out bool
		err error
	}{
		// Only numbergroups
		{"aba", false, ErrNotANumber},
		{"12af", false, ErrNotANumber},
		// Lower bounds
		{"1", false, ErrTooSmall},
		{"999", false, ErrTooSmall},
		// Anything else should be parsable
		{"1000", true, nil},
		{"100000000000", true, nil},
		{"1338054622987", true, nil},
	}

	numbergroup := NewNumberGroup()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := numbergroup.CanParseIntoHuman(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestNumberGroupCanParseFromHuman(t *testing.T) {
	tests := []struct {
		in  string
		out bool
		err error
	}{
		{"999", false, ErrTooSmall},
		{"1000", false, ErrNotHumanGroup},
		{"1,000", true, nil},
		{"1,00f", false, ErrNotANumber},
		{"1,000,000", true, nil},
	}

	numbergroup := NewNumberGroup()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := numbergroup.CanParseFromHuman(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestNumberGroupDoIntoHuman(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"1000", "1,000"},
		{"10000", "10,000"},
		{"100000", "100,000"},
		{"1000000", "1,000,000"},
		{"10000000", "10,000,000"},
		{"100000000", "100,000,000"},
		{"1000000000", "1,000,000,000"},
		{"10000000000", "10,000,000,000"},
		{"100000000000", "100,000,000,000"},
		// This one is to show the logic that's being employed
		{"abcdefghijklmnopqrstuvwxyz", "ab,cde,fgh,ijk,lmn,opq,rst,uvw,xyz"},
	}

	numbergroup := NewNumberGroup()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := numbergroup.DoIntoHuman(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
		})
	}
}

func TestNumberGroupDoFromHuman(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"1,000", "1000"},
		{"10,000", "10000"},
		{"100,000", "100000"},
		{"1,000,000", "1000000"},
		{"10,000,000", "10000000"},
		{"100,000,000", "100000000"},
		{"1,000,000,000", "1000000000"},
		{"10,000,000,000", "10000000000"},
		{"100,000,000,000", "100000000000"},
	}

	numbergroup := NewNumberGroup()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := numbergroup.DoFromHuman(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
		})
	}
}
