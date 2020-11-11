package parsers

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestDefaultUnitIsSet(t *testing.T) {
	tests := []struct {
		in  interface{}
		out string
	}{
		// Should default to 'iec'
		{nil, "iec"},
		// Should be set if allowed
		{"si", "si"},
		{"iec", "iec"},
		// Should ignore bad values and default to iec
		{0123, "iec"},
		{false, "iec"},
		{"notit", "iec"},
	}

	for i, tt := range tests {
		sp := NewSize(tt.in)
		got := sp.units
		if got != tt.out {
			t.Errorf("Case %d: Given = `%v` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
		}
	}
}

func TestSizeCanParseFromMachine(t *testing.T) {
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

	sizeP := NewSize(nil)
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := sizeP.CanParseFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
		})
	}
}

func TestCanParseIntoMachine(t *testing.T) {
	tests := []struct {
		units string
		in    string
		out   bool
	}{
		// Note that the number portion of the input is picked at random in the test
		// table loop. When reading these test cases, you should read them as:
		//
		// 	"$some-gen-number + <in>" -> <out>
		//
		// // @TODO the following isn't happening, it should be added
		// In order to save some redundancy, every upper/lower letter case combination
		// is generated and tested as well. So when we see a letter or pairing
		// of letters like say `KB`, that case will expanded to `kb, KB, kB, Kb`.
		//
		// SI Units
		// Symbols
		{"si", "B", true}, // "B" is not in the SI units, it feels good have inlcude it though
		{"si", "K", true},
		{"si", "M", true},
		{"si", "G", true},
		{"si", "T", true},
		{"si", "P", true},
		{"si", "E", true},
		{"si", "Z", true},
		{"si", "Y", true},
		// "Common" 2 letter symbol
		{"si", "KB", true},
		{"si", "MB", true},
		{"si", "GB", true},
		{"si", "TB", true},
		{"si", "PB", true},
		{"si", "EB", true},
		{"si", "ZB", true},
		{"si", "YB", true},
		// Name
		{"si", "kilo", true},
		{"si", "mega", true},
		{"si", "giga", true},
		{"si", "tera", true},
		{"si", "peta", true},
		{"si", "exa", true},
		{"si", "zetta", true},
		{"si", "yotta", true},
		// IEC Units
		// "Common" 2 letter symbol
		{"iec", "Ki", true},
		{"iec", "Mi", true},
		{"iec", "Gi", true},
		{"iec", "Ti", true},
		{"iec", "Pi", true},
		{"iec", "Ei", true},
		{"iec", "Zi", true},
		{"iec", "Yi", true},
		// Name
		{"iec", "kibi", true},
		{"iec", "mebi", true},
		{"iec", "gibi", true},
		{"iec", "tebi", true},
		{"iec", "pebi", true},
		{"iec", "exbi", true},
		{"iec", "zebi", true},
		{"iec", "yobi", true},
		// Nonsense
		{"iec", "xafadfa", false},
		{"si", "234Af", false},
	}

	for i, tt := range tests {
		sizeP := NewSize(tt.units)
		t.Run(tt.in, func(t *testing.T) {
			num := rand.Intn(10_000_000)
			input := fmt.Sprintf("%d%s", num, tt.in)
			got := sizeP.CanParseIntoMachine(input)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
		})
	}
}

func TestSizeDoFromMachine(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		// Handle "common" sizes (whatever that means lol)
		{"1", "1.0B"},
		{"10", "10.0B"},
		{"100", "100.0B"},
		{"1000", "1.0Kb"},
		{"10000", "10.0Kb"},
		{"100000", "100.0Kb"},
		{"1000000", "1.0Mb"},
		{"10000000", "10.0Mb"},
		{"100000000", "100.0Mb"},
		{"1000000000", "1.0Gb"},
		{"10000000000", "10.0Gb"},
		{"100000000000", "100.0Gb"},
		{"1000000000000", "1.0Tb"},
		{"10000000000000", "10.0Tb"},
		{"100000000000000", "100.0Tb"},
		{"1000000000000000", "1.0Pb"},
		{"10000000000000000", "10.0Pb"},
		{"100000000000000000", "100.0Pb"},
		{"1000000000000000000", "1.0Eb"},
		{"10000000000000000000", "10.0Eb"},
		{"100000000000000000000", "100.0Eb"},
		{"1000000000000000000000", "1.0Zb"},
		{"10000000000000000000000", "10.0Zb"},
		{"100000000000000000000000", "100.0Zb"},
		{"1000000000000000000000000", "1.0Yb"},
		{"10000000000000000000000000", "10.0Yb"},
		{"100000000000000000000000000", "100.0Yb"},
		{"142089140826193550568923157", "142.1Yb"},
	}

	sizeP := NewSize("si")
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := sizeP.DoFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
		})
	}
}

func TestSizeDoFromMachineIEC(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		// Handle "common" sizes (whatever that means lol)
		{"1", "1.0B"},
		{"10", "10.0B"},
		{"100", "100.0B"},
		{"1000", "1.0Ki"},
		{"1024", "1.0Ki"},
		{"10000", "9.8Ki"},
		{"100000", "97.7Ki"},
		{"1000000", "1.0Mi"},
		{"10000000", "9.5Mi"},
		{"100000000", "95.4Mi"},
		{"1000000000", "0.9Gi"},
		{"10000000000", "9.3Gi"},
		{"100000000000", "93.1Gi"},
		{"1000000000000", "0.9Ti"},
		{"10000000000000", "9.1Ti"},
		{"100000000000000", "90.9Ti"},
		{"1000000000000000", "0.9Pi"},
		{"10000000000000000", "8.9Pi"},
		{"100000000000000000", "88.8Pi"},
		{"1000000000000000000", "0.9Ei"},
		{"10000000000000000000", "8.7Ei"},
		{"100000000000000000000", "86.7Ei"},
		{"1000000000000000000000", "0.8Zi"},
		{"10000000000000000000000", "8.5Zi"},
		{"100000000000000000000000", "84.7Zi"},
		{"1000000000000000000000000", "0.8Yi"},
		{"10000000000000000000000000", "8.3Yi"},
		{"100000000000000000000000000", "82.7Yi"},
		{"142089140826193550568923157", "117.5Yi"},
	}

	sizeP := NewSize("iec")
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := sizeP.DoFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
		})
	}
}

// This series is a bit different given that `NewSize` determines the input
// type or standard to use based on the suffix that is passed to it.
// These aren't checking that logic however (other tests do that)
// hence the passing of the `inputType` parameter in the test cases
func TestSizeDoIntoMachine(t *testing.T) {
	tests := []struct {
		inputType interface{}
		in        string
		out       string
	}{
		// IEC
		// bytes
		{"iec", "1b", "1"},
		{"iec", "1B", "1"},
		{"iec", "100B", "100"},
		{"iec", "1000B", "1000"},
		{"iec", "1000000000B", "1000000000"},
		// @TODO add these
		// {"iec", "1Byte", "1"},
		// {"iec", "1BYTE", "1"},
		{"iec", "1k", "1024"},
		{"iec", "10k", "10240"},
		{"iec", "10000k", "10240000"},
		{"si", "10000k", "10000000"},
		{"si", "10000gb", "10000000000000"},
	}

	for i, tt := range tests {
		sizeP := NewSize(tt.inputType)
		t.Run(tt.in, func(t *testing.T) {
			got := sizeP.DoIntoMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
		})
	}
}
