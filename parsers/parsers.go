package parsers

// Parser is the contract that the main command line application will use
type Parser interface {
	CanParseIntoHuman(string) bool
	CanParseFromHuman(string) bool
	DoIntoHuman(string) string
	DoFromHuman(string) string
}
