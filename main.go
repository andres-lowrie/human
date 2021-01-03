package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andres-lowrie/human/format"
	"github.com/andres-lowrie/human/io"
	"github.com/davecgh/go-spew/spew"
)

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

	// Figure out direction and which format
	// we'll default to the `--from` direction since it might be the most common
	// usecase i.e. we want to go "from" machine into human format
	if len(args.Positionals) < 1 {
		log.Warn("@TODO read arguments from stdin")
		return
	}

	input := args.Positionals[0]
	direction := "from"
	format := ""
	for _, d := range []string{"into", "from"} {
		if val, ok := args.Options[d]; ok && val != "" {
			direction = d
			format = val
		}
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
				fmt.Println(output)
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
	if output != "" {
		fmt.Println(output)
	}

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
