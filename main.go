package main

import (
	"fmt"
	"os"

	"github.com/andres-lowrie/human/cmd"
)

func main() {
	// Figure out what was passed into the program
	args := cmd.ParseCliArgs(os.Args[1:])

	// The idea here is that human will print out all parseable values for each
	// known parser (the below map); ie: arguments are used to make it more
	// specific similar to `dig`, where `dig` with no args gives all the
	// information it has, and then something like `dig +short` gives you a whole
	// lot less
	handlers := map[string]cmd.Command{"number": cmd.NewNumber()}

	// Figure out direction and which format
	// we'll default to the `--from` direction since it might be the most common
	// usecase i.e. we want to go "from" machine into human format
	if len(args.Positionals) < 1 {
		fmt.Println("@TODO read arguments from stdin")
		return
	}

	input := args.Positionals[0]
	direction := "from"
	format := ""
	for _, d := range []string{"into", "from"} {
		if val, ok := args.Options[d]; ok && val != "" {
			fmt.Println("val", val)

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
	fmt.Println("format", format)
	fmt.Println("input", input)
	fmt.Println("direction", direction)

	var output string
	if format == "" {
		for _, c := range handlers {
			output = c.Run(direction, input, args)
			if output != "" {
				fmt.Println("Output", output)
			}
		}
		return
	}

	c, ok := handlers[format]
	if !ok {
		fmt.Println("unknown format, nothing to do")
		return
	}

	output = c.Run(direction, input, args)
	if output != "" {
		fmt.Println("Output", output)
	}

}
