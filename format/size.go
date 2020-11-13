package format

import (
	"github.com/andres-lowrie/human/cmd"
	"github.com/andres-lowrie/human/parsers"
)

type Size struct{}

func NewSize() Format {
	return &Size{}
}

func (s *Size) GetParsers() []parsers.Parser {
	// Given that when this method gets called we don't know what units the user
	// is wanting to use, we give back all options so that the parsing check can
	// attempt each of them to see if one will work
	return []parsers.Parser{parsers.NewSize("iec"), parsers.NewSize("si")}
}

func (s *Size) Run(direction, input string, args cmd.CliArgs) (string, error) {
	// We know from the implementation that `iec` is the default so we'll only
	// check for others and default to `iec` if we find nothing
	var p parsers.Parser
	units := args.Options["units"]

	switch units {
	case "si":
		p = parsers.NewSize("si")
	default:
		p = parsers.NewSize("iec")
	}

	if ok, _ := p.CanParseFromMachine(input); direction == "from" && ok {
		return p.DoFromMachine(input), nil
	}

	if ok, _ := p.CanParseIntoMachine(input); direction == "into" && ok {
		return p.DoIntoMachine(input), nil
	}

	return "Err: Input unparsable for `size`", nil
}
