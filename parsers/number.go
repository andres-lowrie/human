package parsers

import (
	"regexp"
	"strings"
)

// NumberGroup handles strings made up of contiguous "0-9" characters converts
// to and from groupings
type NumberGroup struct{}

// NewNumberGroup constructs a NumberGroup struct
func NewNumberGroup() *NumberGroup {
	return &NumberGroup{}
}

// CanParseIntoHuman determines if input is within bounds
// in that the input:
// Contains only digits
// Is >= 1000
func (n *NumberGroup) CanParseIntoHuman(s string) (bool, error) {
	match, _ := regexp.MatchString(`^[0-9]$|^[1-9][0-9]+$`, s)
	if match && len(s) < 4 {
		return false, ErrTooSmall
	} else if match {
		return true, nil
	}

	return false, ErrNotANumber

}

// CanParseFromHuman determines if input is within bounds
// in that the input:
// 	Should at least be the number 1 thousand
// 	Can't have letters in it
// 	Needs to have a comma in it
func (n *NumberGroup) CanParseFromHuman(s string) bool {

	if len(s) <= 4 {
		return false
	}

	if match, _ := regexp.MatchString(`[a-z]+`, s); match {
		return false
	}

	if !strings.Contains(s, ",") {
		return false
	}

	return true
}

// DoIntoHuman takes a string made up of contiguous "0-9" characters and
// returns number groupings
func (n *NumberGroup) DoIntoHuman(s string) string {

	// Figure out where to place commas
	var buf strings.Builder
	bufLen := len(s) - 1

	for i := 0; i <= bufLen; i++ {
		if i%3 == 0 && i != 0 {
			buf.WriteRune(',')
		}
		buf.WriteByte(s[bufLen-i])
	}

	// Output the string
	var out strings.Builder
	outStr := buf.String()
	outLen := len(outStr) - 1

	for i := outLen; i >= 0; i-- {
		out.WriteByte(outStr[i])
	}

	return out.String()
}

// DoFromHuman ...
func (n *NumberGroup) DoFromHuman(s string) string {
	return strings.Replace(s, ",", "", -1)
}
