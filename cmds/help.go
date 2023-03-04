package cmds

import (
	"fmt"
	"strings"

	"github.com/andres-lowrie/human/io"
)

func allCommandOut() []string {
	allCommands := GetAllCommands()

	var rtn []string
	for _, c := range allCommands {
		rtn = append(rtn, fmt.Sprintf("%s: %s", c.Name(), c.ShortDesc()))
	}
	return rtn
}

func allFormatOut() []string {
	allFormats := GetAllFormats()

	var rtn []string
	for _, c := range allFormats {
		rtn = append(rtn, fmt.Sprintf("%s: %s", c.Name(), c.ShortDesc()))
	}
	return rtn
}

type GlobalHelp struct {
	name      string
	usage     string
	shortDesc string
	longDesc  string
}

func NewGlobalHelp() *GlobalHelp {

	return &GlobalHelp{
		name:      "global help",
		usage:     "human [(CMD|FORMAT)] [OPTS] [ARGS...]",
		shortDesc: "human is a utility to translate things from machine formats to human friendly formats.",
		longDesc: fmt.Sprintf(`
CMD can be any of the following
  %s

FORMAT can be any of the following
  %s

OPTS
  Are the flags (short options) and parameters (long options) that the CMD or FORMAT wants.
  There are also the following global options:
    --into: You're supplying human and want to translate into machine
    --from: You're supplying machine and want to translate into human
    -v[vvv]: Verbose output; each v increases the verbosity (from info to debug)
  
ARGS
  This can be positional arguments, path to file(s), or standard in.
  Every FORMAT/CMD defines what it expects,
  for more information run "human help"

Examples
  Add decimal separators to a number:
    human number -g 1000000000

  Translate cron:
    human cron '* * * * *'

  Process every line in a file:
    echo "1024\n4579\n45649" | human size

  Run the input against every format human knows of:
    human $anything
`, strings.Join(allCommandOut(), "\n  "), strings.Join(allFormatOut(), "\n  ")),
	}
}

func (g *GlobalHelp) Name() string {
	return g.name
}

func (g *GlobalHelp) Usage() string {
	return g.usage
}

func (g *GlobalHelp) ShortDesc() string {
	return g.shortDesc
}

func (g *GlobalHelp) LongDesc() string {
	return g.longDesc
}

func (g *GlobalHelp) Run(direction, input string, args io.CliArgs) (string, error) {
	out := UsageTemplate(g)
	return out.String(), nil
}

type Help struct {
	name      string
	usage     string
	shortDesc string
	longDesc  string
	// these are particular for this command
	topics map[string]string
}

func NewHelp() *Help {
	return &Help{
		name:      "help",
		usage:     "human help (CMD|FORMAT|TOPIC)",
		shortDesc: "Get usage and examples formats, commands, and topics",
		longDesc: `
human has 2 types actions arguments, a FORMAT and a CMD. A FORMAT is what
the human tool can translate into and from and a CMD is any action that it
can perform that is not a translation.

Get more in depth help for a CMD|FORMAT running any of these:
  %s
`,
	}
}

func (h *Help) Name() string {
	return h.name
}

func (h *Help) Usage() string {
	return h.usage
}

func (h *Help) ShortDesc() string {
	return h.shortDesc
}

func (h *Help) LongDesc() string {
	var cmds []string
	for _, c := range GetAllCommands() {
		if c.Name() != "help" {
			cmds = append(cmds, c.Name())
		}
	}
	return fmt.Sprintf(h.longDesc, strings.Join(cmds, "\n  "))
}

func (h *Help) Run(direction, input string, args io.CliArgs) (string, error) {
	out := UsageTemplate(h)
	return out.String(), nil
}
