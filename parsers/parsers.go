package parsers

import (
	"errors"
	"regexp"
)

// Parser is the contract that the main command line application will use
type Parser interface {
	CanParseIntoMachine(string) (bool, error)
	CanParseFromMachine(string) (bool, error)
	DoIntoMachine(string) string
	DoFromMachine(string) string
}

// General purpose errors
var ErrNotANumber error = errors.New("Not a Number")
var ErrTooLarge error = errors.New("Too Beaucoup")
var ErrTooSmall error = errors.New("Number too small")

// Catch all error
var ErrUnparsable error = errors.New("Unparsable")

// Empty can be used to as a placeholder for when an
// interface is needed
func NewEmpty() *Empty {
	return &Empty{}
}

type Empty struct{}

func (e *Empty) CanParseIntoMachine(string) (bool, error) {
	return true, nil
}

func (e *Empty) CanParseFromMachine(string) (bool, error) {
	return true, nil
}

func (e *Empty) DoIntoMachine(string) string {
	return "Not Yet Implemented"
}

func (e *Empty) DoFromMachine(string) string {
	return "Not Yet Implemented"
}

// IsMachineNumber validates that a number:
// Contains only digits
// The first digit must be 1-9 when >= 10
func isMachineNumber(s string) bool {
	match, _ := regexp.MatchString(`^[0-9]$|^[1-9][0-9]+$`, s)
	return match
}

// IsDelimitednumber validates that a number:
// Contains only digits and delimiters [., _]
// The first digit must be 1-9 when >= 10
// The first group can be 1-3 digits
// Subsequent groups must be 3 digits
// A number less than 1000 is considered a valid delimited number
func isDelimitedNumber(s string) bool {
	// 0-999 are machine numbers
	if len(s) < 4 {
		if ok := isMachineNumber(s); ok {
			return true
		} else {
			return false
		}
	}

	match, _ := regexp.MatchString(`^[1-9][0-9]{0,2}([.,_ ][0-9]{3})+$`, s)
	return match
}
