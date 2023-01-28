package cmds

import (
	"fmt"
	"strings"
)

type GlobalHelp struct {
	name      string
	usage     string
	shortDesc string
	longDesc  string
}

func NewGlobalHelp() *GlobalHelp {
	allCommands := GetAllCommands()
	allFormats := GetAllFormats()

	var allCommandOut []string
	for _, c := range allCommands {
		allCommandOut = append(allCommandOut, fmt.Sprintf("%s: %s", c.Name(), c.ShortDesc()))
	}

	var allFormatOut []string
	for _, c := range allFormats {
		allFormatOut = append(allFormatOut, fmt.Sprintf("%s: %s", c.Name(), c.ShortDesc()))
	}

	return &GlobalHelp{
		name:      "global help",
		usage:     "human [(CMD|FORMAT)] [OPTS] [ARGS...]",
		shortDesc: "human is a utility to translate things from machine formats to human friendly formats.",
		longDesc: fmt.Sprintf(`
CMD can be any of the following
  %s

FORMAT can be any of the following
  %s
`, strings.Join(allCommandOut, "\n  "), strings.Join(allFormatOut, "\n  ")),
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

type Help struct {
	name      string
	usage     string
	shortDesc string
	longDesc  string
}

func NewHelp() *Help {
	return &Help{
		name:      "help",
		usage:     "human help (CMD|FORMAT)",
		shortDesc: "Get usage and examples for a format or command",
		longDesc:  "@TODO probably need to figure out how to get this information out as well",
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
	return h.longDesc
}
