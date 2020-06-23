package main

import (
	"fmt"
	"os"

	"github.com/andres-lowrie/human/cmd"
)

func main() {
	// Figure out what was passed into the program
	// @TODO add validation process to args, eg make sure format is something we
	// know, multiple directions aren't passed etc. this way we can reduce the
	// amount of checks later in the code and keep all that logic here
	args := cmd.ParseCliArgs(os.Args[1:])

	// The idea here is that human will print out all parseable values for each
	// known parser (the below map); ie: arguments are used to make it more
	// specific similar to `dig`, where `dig` with no args gives all the
	// information it has, and then something like `dig +short` gives you a whole
	// lot less
	handlers := map[string]cmd.Command{"number": cmd.NewNumber()}

	// Figure out direction and which format
	// we'll default to the `--into` direction since it might be the most common
	// usecase i.e. we want to go "into" human format
	input := args.Positionals[0]
	direction := "into"
	format := ""
	for _, d := range []string{"from", "into"} {
		if val, ok := args.Options[d]; ok && val != "" {
			fmt.Println("val", val)

			direction = d
			format = val
		}
	}

	// The logic here is that if no explicit `into` or `from` option was given
	// then the first positional argument (read left from right) is the format
	// and anything after that is the actual input, however if only 1 positional
	// argument was given then that must be the input in which case we should all
	// the possible translations, this is why we're checking format for emptiness
	// twice
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
	if format != "" {
		output = handlers[format].Run(direction, input, args)
	} else {
		for _, c := range handlers {
			output = c.Run(direction, input, args)
		}
	}

	if output != "" {
		fmt.Println("Output", output)
	}
}
