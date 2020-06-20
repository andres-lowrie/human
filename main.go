package main

import (
	"fmt"
	"os"
)

func main() {
	// Figure out what was passed into the program
	args := ParseCliArgs(os.Args[1:])
	fmt.Println("args", args)

	// // @TODO only do this if something more specific wasn't determine
	// the idea here is that human will print out all parseable values for each
	// knows parser (the below map); ie: arguments are used to make it more
	// specific similar to `dig`, where `dig` with no args gives all the
	// information it has, and then something like `dig +short` gives you a while
	// lot less
	// handlers := map[string]parsers.Parser{"number": parsers.NewNumber()}
	//
	// for _, h := range handlers {
	//   fmt.Println(h)
	// }
}
