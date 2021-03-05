package cli

import (
	"errors"
	"flag"
	"os"
)

// Argument contains parsed cli data
type Argument struct {
	File       string
	OutputPath string
}

// checkPath checks if the file/folder represented by the given path exists
func checkPath(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

// Parse parses the cli arguments
func Parse() (Argument, error) {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		return Argument{}, errors.New("missing arguments")
	}
	if checkPath(args[0]) || checkPath(args[1]) {
		return Argument{}, errors.New("bad arguments")
	}
	result := Argument{args[0], args[1]}
	return result, nil
}
