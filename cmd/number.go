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
	return []parsers.Parser{parsers.NewNumberGroup(), parsers.NewNumberWord()}
}

func (n *Number) Run(direction string, input string, args CliArgs) string {
	// Figure out which of the parsers we're using
	var p parsers.Parser
	if _, ok := args.Flags["g"]; ok {
		p = parsers.NewNumberGroup()
	} else if _, ok := args.Flags["w"]; ok {
		p = parsers.NewNumberWord()
	} else {
		p = parsers.NewEmpty()
	}

	if ok, _ := p.CanParseFromHuman(input); direction == "from" && ok {
		return p.DoFromHuman(input)
	}

	if ok, _ := p.CanParseIntoHuman(input); direction == "into" && ok {
		return p.DoIntoHuman(input)
	}

	return ""
}
