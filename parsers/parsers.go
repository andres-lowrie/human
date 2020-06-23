package parsers

// Parser is the contract that the main command line application will use
type Parser interface {
	CanParseIntoHuman(string) bool
	CanParseFromHuman(string) bool
	DoIntoHuman(string) string
	DoFromHuman(string) string
}

// Empty can be used to as a placeholder for when an
// interface is needed
func NewEmpty() *Empty {
	return &Empty{}
}

type Empty struct{}

func (e *Empty) CanParseIntoHuman(string) bool {
	return true
}

func (e *Empty) CanParseFromHuman(string) bool {
	return true
}

func (e *Empty) DoIntoHuman(string) string {
	return "Not Yet Implemented"
}

func (e *Empty) DoFromHuman(string) string {
	return "Not Yet Implemented"
}
