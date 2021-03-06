package format

import (
	"testing"

	"github.com/andres-lowrie/human/io"
	"github.com/andres-lowrie/human/parsers"
)

func TestNumberFormatRun(t *testing.T) {
	tests := []struct {
		direction string
		input     string
		args      io.CliArgs
		out       string
		err       error
	}{
		// Should default to group
		{"from", "10000", io.ParseCliArgs([]string{""}), "10,000", nil},
		{"into", "100,000,000", io.ParseCliArgs([]string{""}), "100000000", nil},
		// Should return error on bad input
		{"from", "notanumber", io.ParseCliArgs([]string{"-w"}), "", parsers.ErrUnparsable},
		// Should return an error when nonsense input is detected
		{"into", "xxxx", io.ParseCliArgs([]string{"-g"}), "", parsers.ErrUnparsable},
		// Happy path
		{"into", "250 thousand", io.ParseCliArgs([]string{"-w"}), "250000", nil},
		{"from", "250000", io.ParseCliArgs([]string{"-w"}), "250 thousand", nil},
		{"into", "250,000", io.ParseCliArgs([]string{"-g"}), "250000", nil},
		{"from", "250000", io.ParseCliArgs([]string{"-g"}), "250,000", nil},
	}
	number := NewNumber()
	for i, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			got, err := number.Run(tt.direction, tt.input, tt.args)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` Args = `%v+`; want `%s` ; got `%s`", i, tt.input, tt.args, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Error Case %d: Given = `%s` Args = `%v+`; want `%t` ; got `%t`", i, tt.input, tt.args, tt.err, err)
			}
		})
	}
}
