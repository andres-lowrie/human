package parsers

import (
	"errors"
	"strings"
)

var ErrNotHumanGroup error = errors.New("Not a Delimited Number")

// NumberGroup handles strings made up of contiguous "0-9" characters converts
// to and from groupings
type NumberGroup struct{}

// NewNumberGroup constructs a NumberGroup struct
func NewNumberGroup() *NumberGroup {
	return &NumberGroup{}
}

// CanParseFromMachine determines if input is within bounds
// in that the input:
// 	Contains only digits
// 	Is >= 1000
func (n *NumberGroup) CanParseFromMachine(s string) (bool, error) {
	if isMachineNumber(s) && len(s) >= 4 {
		return true, nil
	}

	// Error cases
	var err error

	if isMachineNumber(s) && len(s) < 4 {
		err = ErrTooSmall
	} else {
		err = ErrNotANumber
	}
	return false, err
}

// CanParseIntoMachine determines if input is within bounds
// in that the input:
// 	Should at least be the number 1 thousand
// 	Can't have letters in it
// 	Needs to have a comma in it
func (n *NumberGroup) CanParseIntoMachine(s string) (bool, error) {
	if isDelimitedNumber(s) && len(s) >= 4 {
		return true, nil
	}

	// Error cases
	var err error

	if isMachineNumber(s) && len(s) < 4 {
		err = ErrTooSmall
	} else if isMachineNumber(s) {
		err = ErrNotHumanGroup
	} else {
		err = ErrNotANumber
	}
	return false, err
}

// DoFromMachine takes a string made up of contiguous "0-9" characters and
// returns number groupings
func (n *NumberGroup) DoFromMachine(s string) string {
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

// DoIntoMachine ...
func (n *NumberGroup) DoIntoMachine(s string) string {
	return strings.Replace(s, ",", "", -1)
}
