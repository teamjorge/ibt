package examples

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/teamjorge/ibt"
)

// Use the default testing file if no ibt files were provided
// Provided files can use wildcards, for example ./telemetry/*.ibt
func getExampleFilePattern() string {
	flag.Parse()

	var filePattern string
	if flag.Arg(0) == "" {
		filePattern = ".testing/valid_test_file.ibt"
	} else {
		filePattern = flag.Arg(0)
	}

	return filePattern
}

func ParseExampleStubs() (ibt.StubGroup, error) {
	filePattern := getExampleFilePattern()

	// Find all files for parsing in case it had included a wildcard
	files, err := filepath.Glob(filePattern)
	if err != nil {
		return nil, fmt.Errorf("could not glob the given input files: %v", err)
	}

	// Parse the files into stubs
	stubs, err := ibt.ParseStubs(files...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stubs for %v. error - %v", files, err)
	}

	return stubs, nil
}
