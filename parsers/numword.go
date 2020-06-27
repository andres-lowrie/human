// Given a number in string form
// Return the number in word form of the greatest power (e.g. 300,100,000,000 => 300.1 Billion)
// Lower limit > Thousands (e.g. xxx,000) # Numbers below this are in hundreds and can be expressed numerically
// Upper limit > Centillion (10^303) ref: https://en.wikipedia.org/wiki/Names_of_large_numbers#Standard_dictionary_numbers

package parsers

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

// NumberWord handles strings made of contiguous "0-9" characters
// strings delimited by [,. ] are accepted
// strings are assumed to have no decimal places
// converts to word strings of the greatest power
type NumberWord struct{}

// NewNumberWord construcs a NumberWord struct
func NewNumberWord() *NumberWord {
	return &NumberWord{}
}

// CanParseIntoHuman ...
func (n *NumberWord) CanParseIntoHuman(s string) bool {
	// is it 4 or more characters? (e.g. is it => 1000)
	if len(s) < 4 {
		return false
	}

	// is it a (delimited[,. ]) number?
	match, _ := regexp.MatchString(`^(([0-9]+)|([0-9]{1,3}[., ])+[0-9]{1,3})$`, s)
	if match {
		// is it less than the max?
		//  300000 max places plus delimiters
		if len(s) < 400000 {
			return true
		}
	}
	// everything else is not a number
	return false
}

// CanParseFromHuman ...
func (n *NumberWord) CanParseFromHuman(s string) bool {
	// is it a digit word combo? (e.g. 48 billion) /\d+ [a-zA-Z]+)
	// is the word in the translation map?
	return false
}

// DoIntoHuman ...
func (n *NumberWord) DoIntoHuman(s string) string {
	trans := map[int]struct {
		name   string
		powers int
	}{
		1: {"hundred", 2}, // not used
		2: {"thousand", 3},
		3: {"million", 6},
		4: {"billion", 9},
		5: {"trillion", 12},
	}

	// Strip delimiters
	r := regexp.MustCompile("[^0-9]")
	s = r.ReplaceAllString(s, "")

	// Use NumberGroup to make an array
	numgroup := NewNumberGroup()
	numbers := strings.Split(numgroup.DoIntoHuman(s), ",")

	var out strings.Builder
	out.WriteString(numbers[0])

	// Round second group to nearest hundreds (i.e. 1 decimal place)
	x, _ := strconv.ParseFloat(numbers[1], 64)
	decimal := int(math.Round(x / 100.0))

	if decimal > 0 {
		out.WriteString("." + strconv.Itoa(decimal))
	}
	out.WriteString(" " + trans[len(numbers)].name)

	return out.String()
}

// DoFromHuman ...
// Only works with highest power (e.g. 100.3 Billion, not 100,300 Million)
func (n *NumberWord) DoFromHuman(s string) string {
	// Split numbers from word
	// Get powers from translation map
	// Return ( numbers * 10^foo )
	return "100"
}
