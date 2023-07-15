package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andres-lowrie/human/cmds"
	"github.com/andres-lowrie/human/format"
	"github.com/andres-lowrie/human/io"
	"github.com/davecgh/go-spew/spew"
)

type Direction struct {
	From bool
	Into bool
	// @todo fix this hack by passing in the Direction type to the other functions
	toStr func() string
}

func getDirection(args io.CliArgs) Direction {
	d := Direction{From: true}

	if args.Flags["i"] || len(args.Options["--input"]) > 0 {
		d.From = false
		d.Into = true
	}

	d.toStr = func() string {
		if d.Into {
			return "into"
		}
		return "from"
	}

	return d
}

// centralize all the writing to stdout here
func doOutput(a interface{}) {
	log := io.NewLogger(io.OFF, false)
  log.Debug("Output:")
	log.Debug(spew.Sdump(a))
  fmt.Println(a)
}

func run(log io.Ourlog, args io.CliArgs) {
	log.Debug("Program start")
	log.Debug(spew.Sdump(args))

	// The idea here is that human will print out all parseable values for each
	// known parser (the below map); ie: arguments are used to make it more
	// specific similar to `dig`, where `dig` with no args gives all the
	// information it has, and then something like `dig +short` gives you a whole
	// lot less
	// @TODO see if we can use GetParsers instead of instantiating directly
	handlers := map[string]format.Format{"number": format.NewNumber(), "size": format.NewSize()}

	// If user is asking for help, then all other input gathering logic is
	// avoided since this would not be a normal run.
	// Note there's also a "help" command, we're not considering that here because
	// that's a command and so it can follow the normal flow
	passedHelp := func() bool {
		_, shortOpt := args.Flags["h"]
		_, longOpt := args.Options["help"]
		if shortOpt || longOpt {
			return true
		}
		return false
	}()

	// @TODO read stdin
	if len(args.Positionals) < 1 && passedHelp == false {
		fmt.Println("@notimplementedyet read arguments from stdin")
		// if stdio is empty, then we need to show usage
		// r := bufio.NewReader(os.Stdin)
		// _, err := r.Peek(10)
		// // spew.Dump(got)
		// spew.Dump(err)
		return
	}

	if passedHelp {
		helpCmd := cmds.NewGlobalHelp()
		output, _ := helpCmd.Run("", "", args)
		doOutput(output)
		return
	}

	// Are we processing a command or a format?
	//
	// Commands don't require parsers or the like so for those we just need to
	// execute them.
	//
	// If we are processing a format, then we need to figure out direction and
	// which format. We'll default to the `--from` direction since it might be
	// the most common usecase i.e. we want to go "from" machine into human
	// format
	input := args.Positionals[0]
	direction := getDirection(args)
	log.Debug("Direction: %s", direction.toStr())

	if c, ok := cmds.GetCommand(input); ok {
		output, err := c.Run(direction.toStr(), input, args)
		if err != nil {
			return
		}
		doOutput(output)
		return
	}

	// The logic here is that if no explicit `into` or `from` flag (ie: `direction`) was given
	// then the first positional argument (read left from right) is the format
	// and anything after that is the actual input, however if only 1 positional
	// argument was given then that must be the input in which case we should run
	// all the possible translations (ie: formats)
	var format string
	if len(args.Positionals) > 1 {
		format = args.Positionals[0]
		input = args.Positionals[1]
	}
	log.Debug("Format: %s", format)

	var output string
	if format == "" {
		for _, c := range handlers {
			output, _ = c.Run(direction.toStr(), input, args)
			if output != "" {
				doOutput(output)
			}
		}
		return
	}

	c, ok := handlers[format]
	if !ok {
		log.Info("unknown format '%s', nothing to do", format)
		return
	}

	output, _ = c.Run(direction.toStr(), input, args)
	doOutput(output)
}

func main() {
	args := io.ParseCliArgs(os.Args[1:])
	log := io.NewLogger(io.OFF, false)

	// Figure out if we have to enable the logger
	if args.Flags["v"] {
		log = io.NewLogger(io.INFO, true)
	}

	if len(args.Options["v"]) > 0 {
		if n, err := strconv.Atoi(args.Options["v"]); err == nil {
			switch n {
			case 2:
				log = io.NewLogger(io.WARN, true)
			default: // allows users to spam the v's
				log = io.NewLogger(io.DEBUG, true)
			}
		}
	}

	run(log, args)
}
