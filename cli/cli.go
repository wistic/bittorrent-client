package cli

import (
	"errors"
	"flag"
	"os"
)

// Argument contains parsed cli data
type Argument struct {
	Torrent string
	Output  string
}

// checkPath checks if the file/folder represented by the given path exists
func checkPath(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

// Parse parses the cli arguments
func Parse() (*Argument, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	output := flag.String("o", workingDir, "Output directory")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return nil, errors.New("Missing arguments")
	}
	torrent := args[0]
	if checkPath(torrent) || checkPath(*output) {
		return nil, errors.New("Bad arguments")
	}
	return &Argument{torrent, *output}, nil
}
