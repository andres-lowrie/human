package main

import (
	"fmt"
	"os"

	"github.com/andres-lowrie/human/cmd"
)

func main() {
	// Figure out what was passed into the program
	args := cmd.ParseCliArgs(os.Args[1:])

	// // @TODO only do this if something more specific wasn't determine
	// the idea here is that human will print out all parseable values for each
	// knows parser (the below map); ie: arguments are used to make it more
	// specific similar to `dig`, where `dig` with no args gives all the
	// information it has, and then something like `dig +short` gives you a whole
	// lot less
	handlers := map[string]cmd.Command{"number": cmd.NewNumber()}

	// figure out direction and which format
	// we'll default to `--from` since that seems like the most common usecase
	// i.e. we want human string back
	fmt.Println("args", args)

	direction := "from"
	format := ""
	for _, d := range []string{"into", "from"} {
		if args.Options[d] != "" {
			direction = d
			format = args.Options[d]
		}
	}

	fmt.Println("direction", direction)
	fmt.Println("format", format)

	// for _, h := range handlers {
	//   for _, p := range h.GetParsers(){
	//   }
	//   fmt.Println(h)
	// }
	fmt.Println("handlers", handlers)

}
