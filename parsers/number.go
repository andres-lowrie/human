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

// CanParseIntoHuman ...
func (n *NumberGroup) CanParseIntoHuman(s string) bool {
	if len(s) < 4 {
		return false
	}

	match, _ := regexp.MatchString(`[a-z]+`, s)
	return !match
}

// CanParseFromHuman ...
func (n *NumberGroup) CanParseFromHuman(s string) bool {

	conds := []func(string) bool{
		// Should at least be the number 1 thousand
		func(s string) bool {
			return len(s) >= 4
		},
		// Can't have letters in it
		func(s string) bool {
			match, _ := regexp.MatchString(`[a-z]+`, s)
			return !match
		},
		// Needs to have a comma in it
		func(s string) bool {
			return strings.Contains(s, ",")
		},
	}

	for _, c := range conds {
		if !c(s) {
			return false
		}
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
