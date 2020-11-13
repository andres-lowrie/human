package format

import (
	"github.com/andres-lowrie/human/cmd"
	"github.com/andres-lowrie/human/parsers"
)

type Number struct {
}

func NewNumber() Format {
	return &Number{}
}

func (n *Number) GetParsers() []parsers.Parser {
	return []parsers.Parser{parsers.NewNumberGroup(), parsers.NewNumberWord()}
}

func (n *Number) Run(direction string, input string, args cmd.CliArgs) (string, error) {
	// Figure out which of the parsers we're using, default to "groupping, -g"
	var p parsers.Parser

	p = parsers.NewNumberGroup()

	if _, ok := args.Flags["w"]; ok {
		p = parsers.NewNumberWord()
	}

	if ok, _ := p.CanParseFromMachine(input); direction == "from" && ok {
		return p.DoFromMachine(input), nil
	}

	if ok, _ := p.CanParseIntoMachine(input); direction == "into" && ok {
		return p.DoIntoMachine(input), nil
	}

	return "", ErrUnparseableInput
}
