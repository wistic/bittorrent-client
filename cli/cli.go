package cli

import (
	"errors"
	"flag"
)

// Argument contains parsed cli data
type Argument struct {
	File string
}

// Parse parses the cli arguments
func Parse() (Argument, error) {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		return Argument{}, errors.New("No '.torrent' file specified")
	}
	result := Argument{
		File: args[0],
	}
	return result, nil
}
