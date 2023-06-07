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

// centralize all the writing to stdout here
func doOutput(a interface{}) {
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
	// Note there's also a "help" command, were not considering that here because
	// since that's a command it can follow the normal flow
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
		fmt.Println("@TODO read arguments from stdin")
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
	// which format we'll default to the `--from` direction since it might be
	// the most common usecase i.e. we want to go "from" machine into human
	// format
	input := args.Positionals[0]
	direction := "from"
	format := ""
	fmt.Println("1")
	spew.Dump(args)
	for _, d := range []string{"into", "from"} {
		fmt.Println("2")
		fmt.Println(d)
		spew.Dump(args.Options[d])
    fmt.Println("=======")

    // @leftoff according to my notes, direction should be the first Option passed and it should default to `--from`
    // so this is not a bug bug instead a misunderstanding from my part
    //
    // okay so what I have to do is:
    //
    // - create a test case that encodes this logic, which means:
    // -  `/human --into number "1 million"` -> 1000000
    // -  `/human --from number "1 million"` -> Error
    // -  `/human number "1 million"` -> Error
		if val, ok := args.Options[d]; ok && val != "" {
			direction = d
			format = val
		}
	}
	spew.Dump(direction)

	if c, ok := cmds.GetCommand(input); ok {
		output, err := c.Run(direction, input, args)
		if err != nil {
			spew.Dump(err)
			return
		}
		doOutput(output)
		return

	}

	// The logic here is that if no explicit `into` or `from` option was given
	// then the first positional argument (read left from right) is the format
	// and anything after that is the actual input, however if only 1 positional
	// argument was given then that must be the input in which case we should run
	// all the possible translations, this is why we're checking format for
	// emptiness twice
	if format == "" {
		if len(args.Positionals) > 1 {
			format = args.Positionals[0]
			input = args.Positionals[1]
		}
	}
	log.Info("format is set to: ", format)
	log.Info("input is set to: ", input)
	log.Info("direction is set to: ", direction)

	var output string
	if format == "" {
		for _, c := range handlers {
			output, _ = c.Run(direction, input, args)
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

	output, _ = c.Run(direction, input, args)
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
