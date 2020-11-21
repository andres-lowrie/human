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
		err error
	}{
		// At least a thousand
		{"100", false, ErrTooSmall},
		// Only numbers
		{"afsdafa", false, ErrNotANumber},
		{"20M", false, ErrNotANumber},
		// Case insensitive
		{"20m", false, ErrNotANumber},
		// Positive
		{"1000", true, nil},
		{"2048", true, nil},
		{"20484046", true, nil},
	}

	sizeP := NewSize(nil)
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := sizeP.CanParseFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; wanted error `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestCanParseIntoMachineHappyPath(t *testing.T) {
	tests := []struct {
		units string
		in    string
		out   bool
		err   error
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
		{"si", "B", true, nil}, // "B" is not in the SI units, it feels good have inlcude it though
		{"si", "K", true, nil},
		{"si", "M", true, nil},
		{"si", "G", true, nil},
		{"si", "T", true, nil},
		{"si", "P", true, nil},
		{"si", "E", true, nil},
		{"si", "Z", true, nil},
		{"si", "Y", true, nil},
		// "Common" 2 letter symbol
		{"si", "KB", true, nil},
		{"si", "MB", true, nil},
		{"si", "GB", true, nil},
		{"si", "TB", true, nil},
		{"si", "PB", true, nil},
		{"si", "EB", true, nil},
		{"si", "ZB", true, nil},
		{"si", "YB", true, nil},
		// Name
		{"si", "kilo", true, nil},
		{"si", "mega", true, nil},
		{"si", "giga", true, nil},
		{"si", "tera", true, nil},
		{"si", "peta", true, nil},
		{"si", "exa", true, nil},
		{"si", "zetta", true, nil},
		{"si", "yotta", true, nil},
		// IEC Units
		// "Common" 2 letter symbol
		{"iec", "Ki", true, nil},
		{"iec", "Mi", true, nil},
		{"iec", "Gi", true, nil},
		{"iec", "Ti", true, nil},
		{"iec", "Pi", true, nil},
		{"iec", "Ei", true, nil},
		{"iec", "Zi", true, nil},
		{"iec", "Yi", true, nil},
		// Name
		{"iec", "kibi", true, nil},
		{"iec", "mebi", true, nil},
		{"iec", "gibi", true, nil},
		{"iec", "tebi", true, nil},
		{"iec", "pebi", true, nil},
		{"iec", "exbi", true, nil},
		{"iec", "zebi", true, nil},
		{"iec", "yobi", true, nil},
		// It should allow decimal in the number
		{"iec", "100.50ki", true, nil},
	}

	for i, tt := range tests {
		sizeP := NewSize(tt.units)
		t.Run(tt.in, func(t *testing.T) {
			num := rand.Intn(10_000_000)
			input := fmt.Sprintf("%d%s", num, tt.in)

			got, err := sizeP.CanParseIntoMachine(input)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestCanParseIntoMachineEdgeCases(t *testing.T) {
	tests := []struct {
		units string
		in    string
		out   bool
		err   error
	}{
		// Nonsense
		{"iec", "xafadfa", false, ErrUnparsable},
		{"si", "234Af", false, ErrUnknownSuffix},
		// Ensure regex pulls format correctly
		{"iec", "abv0ki", false, ErrUnparsable},
		// It should only allow 1 decimal place
		{"iec", "100.50.3ki", false, ErrUnparsable},
	}

	for i, tt := range tests {
		sizeP := NewSize(tt.units)
		t.Run(tt.in, func(t *testing.T) {
			got, err := sizeP.CanParseIntoMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestSizeDoFromMachine(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err error
	}{
		// Handle "common" sizes (whatever that means lol)
		{"1", "1.0B", nil},
		{"10", "10.0B", nil},
		{"100", "100.0B", nil},
		{"1000", "1.0Kb", nil},
		{"10000", "10.0Kb", nil},
		{"100000", "100.0Kb", nil},
		{"1000000", "1.0Mb", nil},
		{"10000000", "10.0Mb", nil},
		{"100000000", "100.0Mb", nil},
		{"1000000000", "1.0Gb", nil},
		{"10000000000", "10.0Gb", nil},
		{"100000000000", "100.0Gb", nil},
		{"1000000000000", "1.0Tb", nil},
		{"10000000000000", "10.0Tb", nil},
		{"100000000000000", "100.0Tb", nil},
		{"1000000000000000", "1.0Pb", nil},
		{"10000000000000000", "10.0Pb", nil},
		{"100000000000000000", "100.0Pb", nil},
		{"1000000000000000000", "1.0Eb", nil},
		{"10000000000000000000", "10.0Eb", nil},
		{"100000000000000000000", "100.0Eb", nil},
		{"1000000000000000000000", "1.0Zb", nil},
		{"10000000000000000000000", "10.0Zb", nil},
		{"100000000000000000000000", "100.0Zb", nil},
		{"1000000000000000000000000", "1.0Yb", nil},
		{"10000000000000000000000000", "10.0Yb", nil},
		{"100000000000000000000000000", "100.0Yb", nil},
		{"142089140826193550568923157", "142.1Yb", nil},
	}

	sizeP := NewSize("si")
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := sizeP.DoFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}

func TestSizeDoFromMachineIEC(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err error
	}{
		// Handle "common" sizes (whatever that means lol)
		{"1", "1.0B", nil},
		{"10", "10.0B", nil},
		{"100", "100.0B", nil},
		{"1000", "1.0Ki", nil},
		{"1024", "1.0Ki", nil},
		{"10000", "9.8Ki", nil},
		{"100000", "97.7Ki", nil},
		{"1000000", "1.0Mi", nil},
		{"10000000", "9.5Mi", nil},
		{"100000000", "95.4Mi", nil},
		{"1000000000", "0.9Gi", nil},
		{"10000000000", "9.3Gi", nil},
		{"100000000000", "93.1Gi", nil},
		{"1000000000000", "0.9Ti", nil},
		{"10000000000000", "9.1Ti", nil},
		{"100000000000000", "90.9Ti", nil},
		{"1000000000000000", "0.9Pi", nil},
		{"10000000000000000", "8.9Pi", nil},
		{"100000000000000000", "88.8Pi", nil},
		{"1000000000000000000", "0.9Ei", nil},
		{"10000000000000000000", "8.7Ei", nil},
		{"100000000000000000000", "86.7Ei", nil},
		{"1000000000000000000000", "0.8Zi", nil},
		{"10000000000000000000000", "8.5Zi", nil},
		{"100000000000000000000000", "84.7Zi", nil},
		{"1000000000000000000000000", "0.8Yi", nil},
		{"10000000000000000000000000", "8.3Yi", nil},
		{"100000000000000000000000000", "82.7Yi", nil},
		{"142089140826193550568923157", "117.5Yi", nil},
	}

	sizeP := NewSize("iec")
	for i, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := sizeP.DoFromMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
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
		err       error
	}{
		// bytes
		{"iec", "1b", "1", nil},
		{"iec", "1B", "1", nil},
		{"iec", "100B", "100", nil},
		{"iec", "1000B", "1000", nil},
		{"iec", "1000000000B", "1000000000", nil},
		// @TODO add these
		// {"iec", "1Byte", "1",nil},
		// {"iec", "1BYTE", "1",nil},
		{"iec", "1k", "1024", nil},
		{"iec", "10k", "10240", nil},
		{"iec", "10000k", "10240000", nil},
		{"si", "10000k", "10000000", nil},
		{"si", "10000gb", "10000000000000", nil},
		{"iec", "1x", "", ErrUnknownSuffix},
	}

	for i, tt := range tests {
		sizeP := NewSize(tt.inputType)
		t.Run(tt.in, func(t *testing.T) {
			got, err := sizeP.DoIntoMachine(tt.in)
			if got != tt.out {
				t.Errorf("Case %d: Given = `%s` ; want `%s` ; got `%s`", i, tt.in, tt.out, got)
			}
			if err != tt.err {
				t.Errorf("Case %d: Given = `%s` ; want `%t` ; got `%t`", i, tt.in, tt.err, err)
			}
		})
	}
}
