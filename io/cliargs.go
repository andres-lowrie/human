package io

import (
	"strconv"
	"strings"
)

// CliArgs holds the arguments passed into the program
type CliArgs struct {
	Flags       map[string]bool
	Options     map[string]string
	Positionals []string
}

// NewCliArgs creates a new CliArgs struct with initialized values
func NewCliArgs() CliArgs {
	return CliArgs{
		Flags:       map[string]bool{},
		Options:     map[string]string{},
		Positionals: []string{},
	}
}

// IsFlag deterimnes if human considers the string a flag
func IsFlag(s string) bool {
	return len(strings.Split(s, "-")) == 2
}

// IsOption deterimnes if human considers the string an option
func IsOption(s string) bool {
	return strings.HasPrefix(s, "--")
}

// IsPositional deterimnes if human considers the string to be a
// positional parameter
func IsPositional(s string) bool {
	return !strings.HasPrefix(s, "-")
}

// ParseCliArgs transforms the slice of strings passed into the program into
// the concrete type all the parsers know how to deal with.
func ParseCliArgs(input []string) CliArgs {
	args := NewCliArgs()
	for i := 0; i <= len(input)-1; i++ {
		v := input[i]
		// We're going to consume the input and break out each word into
		// some type defined in the CliArgs struct above.
		//
		// We're going to keep the actual content as primitive types so
		// that the parsers that actually do the work can be loosely
		// coupled to this process (the process of gathering arguments)

		if IsPositional(v) {
			args.Positionals = append(args.Positionals, v)
			continue
		}

		if IsOption(v) {
			var pair []string
			var key string
			var value string

			pair = strings.Split(strings.TrimLeft(v, "--"), "=")
			key = pair[0]
			value = "" // @NOTE We could in the future not deafult to this and instead fail

			if len(pair) == 2 {
				value = pair[1]
			}

			// Should we swallow the next positional?
			if len(pair) != 2 {
				nextWordIdx := i + 1
				if nextWordIdx < len(input) {
					if IsPositional(input[nextWordIdx]) {
						value = input[nextWordIdx]
						i = nextWordIdx
					}
				}
			}

			args.Options[key] = value
		}

		if IsFlag(v) {
			flags := strings.Split(strings.TrimLeft(v, "-"), "")
			for _, v := range flags {
				// repeated flags get summed up and are available as options with their
				// count instead of flags
				if args.Flags[v] == true {
					args.Options[v] = "2"
					delete(args.Flags, v)
				} else if len(args.Options[v]) > 0 {
					n, _ := strconv.Atoi(args.Options[v])
					n++
					args.Options[v] = strconv.Itoa(n)
				} else {
					args.Flags[v] = true
				}

			}
		}
	}
	return args
}
