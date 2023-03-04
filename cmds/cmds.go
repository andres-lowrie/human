package cmds

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/andres-lowrie/human/format"
	"github.com/andres-lowrie/human/io"
)

// ... can't think of an easier way to enforce fields on Commands without
// writing all these "Getters" :(
type Command interface {
	Name() string
	Usage() string
	ShortDesc() string
	LongDesc() string
	Run(string, string, io.CliArgs) (string, error)
}

// GetAllCommands returns all the commands. Note that `Command` is not the same
// as a `Format`; a command is a cli "sub-command" whereas a format is the things
// human knows how to parse
//
// Note that this skips over the GlobalHelp command because that command uses
// this function at initialization
func GetAllCommands() map[string]Command {
	return map[string]Command{"help": NewHelp()}
}

// GetAllFormats returns all the formats as Commands, given
// that they also implement the Command interface. This is
// used to build usage strings
func GetAllFormats() map[string]Command {
	return map[string]Command{"number": format.NewNumber(), "size": format.NewSize()}
}

// GetCommand checks a string against the Name of every command, returns the
// command if they match
func GetCommand(s string) (Command, bool) {
	for _, c := range GetAllCommands() {
		if c.Name() == s {
			return c, true
		}
	}
	return nil, false
}

// UsageTemplate used to create strings for output from cmds.Command
func UsageTemplate(c Command) bytes.Buffer {
	funcs := template.FuncMap{
		"rmLead": func(s string) string { return strings.TrimLeft(s, "\n") },
	}

	tpl := strings.TrimLeft(`
usage: {{.Usage |rmLead}}

{{.ShortDesc|rmLead}}

{{.LongDesc|rmLead}}
`, "\n")

	t := template.Must(template.New("usage").Funcs(funcs).Parse(tpl))
	var b bytes.Buffer
	t.Execute(&b, c)
	return b

}
