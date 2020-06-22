package cmd

import (
	"github.com/andres-lowrie/human/parsers"
)

type Number struct {
}

func NewNumber() Command {
	return &Number{}
}

func (n *Number) GetParsers() []parsers.Parser {
	return []parsers.Parser{parsers.NewNumberGroup()}
}

func (n *Number) Run(direction string, input string, args CliArgs) string {
	// Figure out which of the parsers we're using
	var p parsers.Parser
	if _, ok := args.Flags["g"]; ok {
		p = parsers.NewNumberGroup()
	} else {
		// @TODO implement numberword parser
	}

	if direction == "from" && p.CanParseFromHuman(input) {
		return p.DoFromHuman(input)
	}

	if direction == "into" && p.CanParseIntoHuman(input) {
		return p.DoIntoHuman(input)
	}

	return ""
}
