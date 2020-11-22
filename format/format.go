package format

import (
	"github.com/andres-lowrie/human/io"
	"github.com/andres-lowrie/human/parsers"
)

type Format interface {
	GetParsers() []parsers.Parser
	Run(string, string, io.CliArgs) (string, error)
}
