package format

import (
	"github.com/andres-lowrie/human/io"
	"github.com/andres-lowrie/human/parsers"
	"github.com/davecgh/go-spew/spew"
)

type Size struct {
	name      string
	usage     string
	shortDesc string
	longDesc  string
}

func NewSize() *Size {
	return &Size{
		name:      "size",
		usage:     "human size [--units (si|iec)] [ARGS]...",
		shortDesc: "converts continuous numbers into common machine sizes like Mb Gb or converts sizes into number of bytes they represent",
		longDesc: `
When "--units" isn't passed, it defaults to "iec" (See tables below for more info).

From Machine:
  Typically when dealing with machine output  most things expect the number of bytes.

  Examples:
    How many gibibytes is 1000000000
    """
    echo 1000000000 | human
    > 0.9Gi
    """

    How many gigabytes is 1000000000
    """
    echo 1000000000 | human --unit si
    > 1.0Gb
    """

To Machine:
  This expects a suffix at the end of the input to know how to respond

  Examples:
    Give me the number of bytes in 2Mi
    """
    human 2Mi
    > 2097152
    """

    Give me the number of bytes in 2G using si conversion
    """
    human --units si 2g
    2000000000
    """

Units Breakdown:
  The following tables shows how the units work. The columns "short" and "long"
  show the suffixes that can be given and the column "value" shows the factor of
  the size
  
          units: si                             unit: iec
  | short |  long  | value |           | short |  long  | value  |
  |=======|========|=======|           |=======|========|========|
  |   k   | kilo   | 10^3  |           |  ki   |  kibi  |  2^10  |
  |   m   | mega   | 10^6  |           |  mi   |  mebi  |  2^20  |
  |   g   | giga   | 10^9  |           |  gi   |  gibi  |  2^30  |
  |   t   | terra  | 10^12 |           |  ti   |  tebi  |  2^40  |
  |   p   | peta   | 10^15 |           |  pi   |  pebi  |  2^50  |
  |   e   | exa    | 10^18 |           |  ei   |  exbi  |  2^60  |
  |   z   | zetta  | 10^21 |           |  zi   |  zebi  |  2^70  |
  |   y   | yotta  | 10^24 |           |  yi   |  yobi  |  2^80  |
  |   r   | ronna  | 10^27 |           
  |   q   | quetta | 10^30 |           
`,
	}
}

func (s *Size) Name() string {
	return s.name
}

func (s *Size) Usage() string {
	return s.usage
}

func (s *Size) ShortDesc() string {
	return s.shortDesc
}

func (s *Size) LongDesc() string {
	return s.longDesc
}

func (s *Size) GetParsers() []parsers.Parser {
	// Given that when this method gets called we don't know what units the user
	// is wanting to use, we give back all options so that the parsing check can
	// attempt each of them to see if one will work
	return []parsers.Parser{parsers.NewSize("iec"), parsers.NewSize("si")}
}

func (s *Size) Run(direction, input string, args io.CliArgs) (string, error) {
	log := io.GetLogger(args)

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
	log.Debug("Parser: %s", spew.Sdump(p))

	// @todo: change this to only check given the direction we're going
	// so that we can give back the error
	if ok, _ := p.CanParseFromMachine(input); direction == "from" && ok {
		return p.DoFromMachine(input)
	}

	if ok, _ := p.CanParseIntoMachine(input); direction == "into" && ok {
		return p.DoIntoMachine(input)
	}

	return "", parsers.ErrUnparsable
}
