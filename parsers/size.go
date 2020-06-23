package parsers

import (
	"regexp"
)

// Size can convert numbers to and from bytes
type Size struct{}

// NewSize constructs a Size parser
func NewSize() *Size {
	return &Size{}
}

// CanParseIntoHuman determines if string is vaid for this parser
func (sz *Size) CanParseIntoHuman(s string) bool {
	if len(s) < 4 {
		return false
	}

	match, _ := regexp.MatchString(`[a-z]+`, s)
	return !match
}

func (sz *Size) DoIntoHuman(s string) string {
	// In theory we could take log10 of the number and figure out how many
	// multiples of 10 we're dealing with then take that number and then do a
	// lookup anyway to see which suffix to use, but doing the string
	// manipulation seems easier... and less steps... but perhaps I'm too stupid
	// to see it some other way
	var suffix string
	switch length := len(s); {
	case length < 4:
		suffix = "B"
	case length >= 4 && length <= 6:
		suffix = "Kb"
	case length >= 7 && length <= 9:
		suffix = "Mb"
	case length >= 10 && length <= 12:
		suffix = "Gb"
	case length >= 13 && length <= 15:
		suffix = "Tb"
	case length >= 16 && length <= 18:
		suffix = "Pb"
	case length >= 19 && length <= 21:
		suffix = "Eb"
	case length >= 22 && length <= 24:
		suffix = "Zb"
	case length >= 25 && length <= 27:
		suffix = "Yb"
	}

	return suffix
}
