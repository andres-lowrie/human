package parsers

import (
	"testing"
)

func TestNumberGroupCanParseFromMachine(t *testing.T) {
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
			got, err := numbergroup.CanParseFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestNumberGroupCanParseIntoMachine(t *testing.T) {
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
			got, err := numbergroup.CanParseIntoMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestNumberGroupDoFromMachine(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err error
	}{
		{"1000", "1,000", nil},
		{"10000", "10,000", nil},
		{"100000", "100,000", nil},
		{"1000000", "1,000,000", nil},
		{"10000000", "10,000,000", nil},
		{"100000000", "100,000,000", nil},
		{"1000000000", "1,000,000,000", nil},
		{"10000000000", "10,000,000,000", nil},
		{"100000000000", "100,000,000,000", nil},
		// This one is to show the logic that's being employed
		{"abcdefghijklmnopqrstuvwxyz", "ab,cde,fgh,ijk,lmn,opq,rst,uvw,xyz", nil},
	}

	numbergroup := NewNumberGroup()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := numbergroup.DoFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestNumberGroupDoIntoMachine(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err error
	}{
		{"1,000", "1000", nil},
		{"10,000", "10000", nil},
		{"100,000", "100000", nil},
		{"1,000,000", "1000000", nil},
		{"10,000,000", "10000000", nil},
		{"100,000,000", "100000000", nil},
		{"1,000,000,000", "1000000000", nil},
		{"10,000,000,000", "10000000000", nil},
		{"100,000,000,000", "100000000000", nil},
	}

	numbergroup := NewNumberGroup()
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := numbergroup.DoIntoMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}
