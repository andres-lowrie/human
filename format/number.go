package format

import (
	"github.com/andres-lowrie/human/io"
	"github.com/andres-lowrie/human/parsers"
)

type Number struct {
	name      string
	usage     string
	shortDesc string
	longDesc  string
}

func NewNumber() *Number {
	return &Number{
		name:      "number",
		usage:     "human number [-(g|w)] [ARGS...]",
		shortDesc: "makes continuous numbers easier to read or turns english numbers into base 10 numbers; eg. '1 million' gives '1000000'",
    longDesc:  `
There are 2 ways to deal with numbers: as groups or words and they are controlled by the following flags

-g: Groups means digit grouping and a comma (,) is used separate, thousands, ten thousands, and so on
-w: Words means that you want the english representation of a number

If no flag is given, then this defaults to groups (-g)
    `,
	}
}

func (n *Number) Name() string {
	return n.name
}

func (n *Number) Usage() string {
	return n.usage
}

func (n *Number) ShortDesc() string {
	return n.shortDesc
}

func (n *Number) LongDesc() string {
	return n.longDesc
}

func (n *Number) GetParsers() []parsers.Parser {
	return []parsers.Parser{parsers.NewNumberGroup(), parsers.NewNumberWord()}
}

func (n *Number) Run(direction string, input string, args io.CliArgs) (string, error) {
	// Figure out which of the parsers we're using, default to "groupping, -g"
	var p parsers.Parser

	p = parsers.NewNumberGroup()

	if _, ok := args.Flags["w"]; ok {
		p = parsers.NewNumberWord()
	}

	if ok, _ := p.CanParseFromMachine(input); direction == "from" && ok {
		return p.DoFromMachine(input)
	}

	if ok, _ := p.CanParseIntoMachine(input); direction == "into" && ok {
		return p.DoIntoMachine(input)
	}

	return "", parsers.ErrUnparsable
}
