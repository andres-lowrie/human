package parsers

// Number handles strings made up of contiguous "0-9" characters
type Number struct{}

// NewNumber constructs a Number struct
func NewNumber() *Number {
	return &Number{}
}

// CanParseIntoHuman ...
func (n *Number) CanParseIntoHuman(s string) bool {
	return false
}

// CanParseFromHuman determines if the input is something this parser can handle
func (n *Number) CanParseFromHuman(s string) bool {
	return false
}

// DoIntoHuman takes a string made up of contiguous "0-9" characters and
// returns numbers words or a grouped number
func (n *Number) DoIntoHuman(s string) string {
	return ""
}

// DoFromHuman takes a string made up of @TODO "number words" or "number
// abbreviations" and returns a string made up of contiguous "0-9" characters.
func (n *Number) DoFromHuman(s string) string {
	return ""
}
