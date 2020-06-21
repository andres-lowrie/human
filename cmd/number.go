package cmd

import (
	"github.com/andres-lowrie/human/parsers"
)

type Number struct {
}

func NewNumber() Command {
	return &Number{}
}

func (n *Number) GetParsers() []parsers.Parser {
	return []parsers.Parser{parsers.NewNumberGroup()}
}

func (n *Number) Run(args CliArgs) string {
	return ""
}
