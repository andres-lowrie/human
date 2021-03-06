package format

import (
	"testing"

	"github.com/andres-lowrie/human/io"
	"github.com/andres-lowrie/human/parsers"
)

func TestSizeFormatRun(t *testing.T) {
	tests := []struct {
		direction string
		input     string
		args      io.CliArgs
		out       string
		err       error
	}{
		// Should default to iec
		{"from", "1024", io.ParseCliArgs([]string{""}), "1.0Ki", nil},
		{"into", "1Mi", io.ParseCliArgs([]string{""}), "1048576", nil},
		// Should accept a `units` option
		{"from", "1024", io.ParseCliArgs([]string{"--units", "iec"}), "1.0Ki", nil},
		{"from", "1000", io.ParseCliArgs([]string{"--units", "si"}), "1.0Kb", nil},
		// Should fail if input is unparsable
		{"from", "xxxx", io.ParseCliArgs([]string{""}), "", parsers.ErrUnparsable},
		// Happy Path
		{"from", "2097152", io.ParseCliArgs([]string{"--units", "iec"}), "2.0Mi", nil},
		{"into", "1G", io.ParseCliArgs([]string{"--units", "si"}), "1000000000", nil},
	}

	size := NewSize()
	for i, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			got, err := size.Run(tt.direction, tt.input, tt.args)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` Args = `%v+`; want `%s` ; got `%s`", i, tt.input, tt.args, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Error Case %d: Given = `%s` Args = `%v+`; want `%t` ; got `%t`", i, tt.input, tt.args, tt.err, err)
			}
		})
	}
}
