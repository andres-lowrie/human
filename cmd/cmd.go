package cmd

import "github.com/andres-lowrie/human/parsers"

type Command interface {
	GetParsers() []parsers.Parser
	Run(string, string, CliArgs) string
}
