package cmds

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/andres-lowrie/human/format"
)

// ... can't think of an easier way to enforce fields on Commands without
// writing all these "Getters" :(
type Command interface {
	Name() string
	Usage() string
	ShortDesc() string
	LongDesc() string
}

// GetAllCommands returns all the commands. Note that `Command` is not the same
// as a `Format`; a command is cli "sub-command" whereas a format is the things
// human knows how to parse
//
// Note that this skips over the GlobalHelp command because that command uses
// this function at initialization
func GetAllCommands() []Command {
	return []Command{NewHelp()}
}

// GetAllFormats returns all the formats as Commands, given
// that they also implement the Command interface. This is
// used to build usage strings
func GetAllFormats() []Command {
	return []Command{format.NewNumber(), format.NewSize()}
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
