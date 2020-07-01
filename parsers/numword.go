// Given a number in string form
// Return the number in word form of the greatest power (e.g. 300,100,000,000 => 300.1 Billion)
// Lower limit > Thousands (e.g. xxx,000) # Numbers below this are in hundreds and can be expressed numerically
// Upper limit > Vigintillion (10^63)
// ref: https://en.wikipedia.org/wiki/Names_of_large_numbers#Standard_dictionary_numbers

package parsers

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// NumberWord handles strings made of contiguous "0-9" characters
// strings delimited by [,. ] are accepted
// strings are assumed to have no decimal places
// converts to word strings of the greatest power
type NumberWord struct {
	trans map[int]struct {
		name   string
		powers int
	}
}

// NewNumberWord constructs a NumberWord struct
func NewNumberWord() *NumberWord {
	return &NumberWord{
		trans: map[int]struct {
			name   string
			powers int
		}{
			1:  {"hundred", 2}, // not used
			2:  {"thousand", 3},
			3:  {"million", 6},
			4:  {"billion", 9},
			5:  {"trillion", 12},
			6:  {"quadrillion", 15},
			7:  {"quintillion", 18},
			8:  {"sextillion", 21},
			9:  {"septillion", 24},
			10: {"octillion", 27},
			11: {"nonillion", 30},
			12: {"decillion", 33},
			13: {"undecillion", 36},
			14: {"duodecillion", 39},
			15: {"tredecillion", 42},
			16: {"quattuordecillion", 45},
			17: {"quindecillion", 48},
			18: {"sexdecillion", 51},
			19: {"septendecillion", 54},
			20: {"octodecillion", 57},
			21: {"novemdecillion", 60},
			22: {"vigintillion", 63},
		},
	}
}

// CanParseIntoHuman ...
// is it 4 or more characters? (e.g. is it => 1000)
// is it a (delimited[,. ]) number?
// is it less than the max?
//  67 places plus delimiters = 88 char
// everything else is not a number
func (n *NumberWord) CanParseIntoHuman(s string) (bool, error) {
	match, _ := regexp.MatchString(`^(([0-9]+)|([0-9]{1,3}[., ])+[0-9]{1,3})$`, s)
	if match {
		if len(s) >= 88 {
			return false, ErrTooLarge
		} else if len(s) < 4 {
			return false, ErrTooSmall
		} else {
			return true, nil
		}
	}
	return false, ErrNotANumber
}

// CanParseFromHuman ...
// is it a digit word combo? ( <number>[.tenths] <word> )
// is the word in the trans table? (case insensitive)
func (n *NumberWord) CanParseFromHuman(s string) bool {
	match, _ := regexp.MatchString(`^[0-9]+([.][0-9])? [a-zA-Z]+$`, s)
	if match {
		_, word := splitHumanNumberWord(s)

		for _, v := range n.trans {
			if v.name == word {
				return true
			}
		}
	}

	return false
}

// DoIntoHuman ...
// Can accept delimited numbers
// Uses NumberGroup to make an array
// Rounds second group to nearest hundreds (i.e. 1 decimal place)
func (n *NumberWord) DoIntoHuman(s string) string {
	// Strip delimiters
	r := regexp.MustCompile("[^0-9]")
	s = r.ReplaceAllString(s, "")

	numgroup := NewNumberGroup()
	numbers := strings.Split(numgroup.DoIntoHuman(s), ",")

	var out strings.Builder
	out.WriteString(numbers[0])

	x, _ := strconv.ParseFloat(numbers[1], 64)
	decimal := int(math.Round(x / 100.0))

	if decimal > 0 {
		out.WriteString("." + strconv.Itoa(decimal))
	}
	out.WriteString(" " + n.trans[len(numbers)].name)

	return out.String()
}

// DoFromHuman ...
// Only works with highest power
// and first digit (e.g. 100.3 Billion, not 100,300 Million)
// Returns a numeric string e.g. 1 thousand => 1000
func (n *NumberWord) DoFromHuman(s string) string {
	num, word := splitHumanNumberWord(s)
	var power int

	// Get power from translation map
	for _, v := range n.trans {
		if word == v.name {
			power = v.powers
		}
	}

	var out strings.Builder
	if a := strings.Split(num, "."); len(a) > 1 {
		fmt.Fprintf(&out, "%s%s", a[0], a[1])
		power--
	} else {
		out.WriteString(num)
	}
	out.WriteString(strings.Repeat("0", power))
	return out.String()
}

// splitHumanNumberWord takes a digit word pair and returns the individual components
// <digit>[.<tenths>] <word>
func splitHumanNumberWord(s string) (string, string) {
	a := strings.Split(s, " ")
	num, word := a[0], a[1]
	word = strings.ToLower(word)
	return num, word
}
