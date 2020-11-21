package parsers

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var ErrUnknownSuffix error = errors.New("Unknown unit suffix")

// Suffixes describes which symbols and names are allowed for a given unit
// Note that all of these are case insensitive when it comes time to do the
// checking
var suffixes = map[string][]string{
	"si": []string{
		"b", "k", "m", "g", "t", "p", "e", "z", "y", "kb", "mb",
		"gb", "tb", "pb", "eb", "zb", "yb", "kilo", "mega", "giga",
		"tera", "peta", "exa", "zetta", "yotta",
	},
	"iec": []string{
		"b", "ki", "mi", "gi", "ti", "pi", "ei", "zi", "yi", "kibi", "mebi",
		"gibi", "tebi", "pebi", "exbi", "zebi", "yobi",
	},
}

// Size can convert numbers to and from bytes
type Size struct {
	// units refers to the suffixes and which standard is being used, SI or IEC.
	// See the `trans` field for more info
	units string
	// unitSuffix refers to the `i` of `b` that makes up the full suffix when
	// dealing with any size bigger than B
	unitSuffix string
	// base refers to mathematical base, eg: base 10 , base 2, base 8, etc.
	base float64
	// trans holds the translation mapping between the `length` of a number
	// string and the corresponding calculation settings, like how to determine
	// the exponent base, the suffix of the byte etc.
	trans map[int]struct {
		// suffix is one of B, Kb, Mb etc. depending on the units we're using
		// Internation System of Units (SI) or International Electrotechnical
		// Commission (IEC)
		//
		// Used this for the logic:
		// 	https://en.wikipedia.org/wiki/Megabyte
		// 	https://en.wikipedia.org/wiki/Binary_prefix
		// 	https://en.wikipedia.org/wiki/International_System_of_Units
		//
		// Ultimately my idea was that it would function the same as the `numfmt`
		// unix tool
		// 	https://www.gnu.org/software/coreutils/manual/html_node/numfmt-invocation.html
		// 	section `26.2.2 Possible units`
		suffix string
		power  float64
	}
}

// NewSize constructs a Size parser
// We're using the empty interface type here (as a hack?) to allow an optional
// parameter. The values allowed are nil, "si", or "iec".
//
// Defaults to "iec" if anything unknown is passed
func NewSize(input interface{}) *Size {
	var units string

	def := "iec"
	allowed := map[string]bool{"si": true, def: true}

	if s, ok := input.(string); !ok {
		units = def
	} else {
		if _, ok := allowed[s]; ok {
			units = s
		} else {
			units = def
		}
	}

	switch units {
	case "si":
		return &Size{
			units:      units,
			unitSuffix: "b",
			base:       10.0,
			trans: map[int]struct {
				suffix string
				power  float64
			}{
				0:  {"B", 0.0}, //@TODO should this be 1.0?
				1:  {"B", 0.0},
				2:  {"B", 0.0},
				3:  {"B", 0.0},
				4:  {"Kb", 3.0},
				5:  {"Kb", 3.0},
				6:  {"Kb", 3.0},
				7:  {"Mb", 6.0},
				8:  {"Mb", 6.0},
				9:  {"Mb", 6.0},
				10: {"Gb", 9.0},
				11: {"Gb", 9.0},
				12: {"Gb", 9.0},
				13: {"Tb", 12.0},
				14: {"Tb", 12.0},
				15: {"Tb", 12.0},
				16: {"Pb", 15.0},
				17: {"Pb", 15.0},
				18: {"Pb", 15.0},
				19: {"Eb", 18.0},
				20: {"Eb", 18.0},
				21: {"Eb", 18.0},
				22: {"Zb", 21.0},
				23: {"Zb", 21.0},
				24: {"Zb", 21.0},
				25: {"Yb", 24.0},
				26: {"Yb", 24.0},
				27: {"Yb", 24.0},
			},
		}
	case "iec":
		fallthrough
	default:
		return &Size{
			units:      units,
			unitSuffix: "i",
			base:       2.0,
			trans: map[int]struct {
				suffix string
				power  float64
			}{
				0:  {"B", 0.0},
				1:  {"B", 0.0},
				2:  {"B", 0.0},
				3:  {"B", 0.0},
				4:  {"Ki", 10.0},
				5:  {"Ki", 10.0},
				6:  {"Ki", 10.0},
				7:  {"Mi", 20.0},
				8:  {"Mi", 20.0},
				9:  {"Mi", 20.0},
				10: {"Gi", 30.0},
				11: {"Gi", 30.0},
				12: {"Gi", 30.0},
				13: {"Ti", 40.0},
				14: {"Ti", 40.0},
				15: {"Ti", 40.0},
				16: {"Pi", 50.0},
				17: {"Pi", 50.0},
				18: {"Pi", 50.0},
				19: {"Ei", 60.0},
				20: {"Ei", 60.0},
				21: {"Ei", 60.0},
				22: {"Zi", 70.0},
				23: {"Zi", 70.0},
				24: {"Zi", 70.0},
				25: {"Yi", 80.0},
				26: {"Yi", 80.0},
				27: {"Yi", 80.0},
			},
		}
	}
}

// CanParseFromMachine determines if string is valid for this parser
func (sz *Size) CanParseFromMachine(s string) (bool, error) {
	if match, _ := regexp.MatchString(`[a-zA-Z]+`, s); match {
		return false, ErrNotANumber
	}

	if len(s) < 4 {
		return false, ErrTooSmall
	}

	return true, nil
}

// CanParseIntoMachine determines if the user input can be handled by this parser.
// The gist is that it should allow any number followed by a known suffix,
// with no spaces in between the number of the suffix, ie:
// 	1234654<suffix>
func (sz *Size) CanParseIntoMachine(s string) (bool, error) {
	// Get the suffix passed
	_, inputSuffix, err := getInputComponents(s)
	if err != nil {
		return false, err
	}

	// Does this suffix correspond to the units we're using?
	allowed := false
	for _, v := range suffixes[sz.units] {
		if v == strings.ToLower(inputSuffix) {
			allowed = true
			break
		}
	}

	if !allowed {
		return false, ErrUnknownSuffix
	}

	return true, nil
}

func (sz *Size) DoFromMachine(s string) string {

	opts := sz.trans[len(s)]
	n, _ := strconv.ParseFloat(s, 64)

	denominator := math.Pow(sz.base, opts.power)
	res := n / denominator

	return fmt.Sprintf("%.1f%s", res, opts.suffix)
}

func (sz *Size) DoIntoMachine(s string) string {
	// Pull out the number and the suffix from the string
	num, suffix, err := getInputComponents(s)
	if err != nil {
		return err.Error()
	}

	// Figure out which lookup to use.
	// When looking up the suffixes we only care about the first letter since the
	// second one is always an `i` or an `b` depending on the units.
	//
	// Here the outlier is `b` since it's only 1 letter
	suffix = strings.Title(strings.ToLower(suffix))
	if suffix != "B" && len(suffix) == 1 {
		suffix = suffix + sz.unitSuffix
	}

	var lookup struct {
		suffix string
		power  float64
	}
	for _, v := range sz.trans {
		if suffix == v.suffix {
			fmt.Println("found a match")
			lookup = v
			break
		}
	}

	// If we didn't find the suffix then we don't know how to handle it
	if lookup.suffix == "" {
		return ""
	}

	multiplier := math.Pow(sz.base, lookup.power)
	res := num * multiplier

	return fmt.Sprintf("%d", int(res))
}

// getInputComponents splits out the input string into the expected components:
// 	the number
// 	and the size suffix
func getInputComponents(s string) (float64, string, error) {
	r := regexp.MustCompile(`(?i)^([0-9]+)(\.[0-9]+)?([a-z]+)`)
	match := r.FindStringSubmatch(s)

	switch len(match) {
	case 4:
		// we got decimal
		num, err := strconv.ParseFloat(match[1]+match[2], 64)
		suffix := match[3]
		return num, suffix, err
	case 3:
		// no decimal
		num, err := strconv.ParseFloat(match[1], 64)
		suffix := match[2]
		return num, suffix, err
	default:
		return 0, "", ErrUnparsable
	}
}
