package format

import (
	"errors"

	"github.com/andres-lowrie/human/cmd"
	"github.com/andres-lowrie/human/parsers"
)

type Format interface {
	GetParsers() []parsers.Parser
	Run(string, string, cmd.CliArgs) (string, error)
}

var ErrUnparseableInput error = errors.New("Could not parse input")
