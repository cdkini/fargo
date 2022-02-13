package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	args := gatherArgs()
	results := runRipgrep(args)
	for _, result := range results {
		promptUserConfirmation(result.text)
	}
}

type SearchResult struct {
	path string
	line int
	text string
}

// <FILE_PATH>:<LINE_NUMBER>:<MATCHING_TEXT>
var rgx = regexp.MustCompile(`(.+):(\d+):(.+)`)

func parseResultFromString(str string) (SearchResult, error) {
	matches := rgx.FindStringSubmatch(str)
	if len(matches) != 4 {
		err := errors.New("Something went wrong when parsing text")
		return SearchResult{}, err
	}

	path := matches[1]
	line, err := strconv.Atoi(matches[2])
	text := matches[3]
	if err != nil {
		log.Fatal(err) // Should never occur since ripgrep guarantees this format
	}

	return SearchResult{path, line, text}, nil
}

func gatherArgs() []string {
	args := os.Args[1:]
	args = append(args, "-n") // Line numbers are essential for proper SearchResult parsing
	return args
}

func runRipgrep(args []string) []SearchResult {
	rg := exec.Command("rg", args...)
	out, err := rg.Output()
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(out), "\n")

	results := make([]SearchResult, 0)
	for _, line := range lines {
		result, err := parseResultFromString(line)
		if err != nil {
			continue // Failed parsing shouldn't cause a sys exit
		}
		results = append(results, result)
	}

	return results
}

// TODO(cdkini): This needs refinement
func promptUserConfirmation(message string) bool {
	fmt.Println(message)
	var response string
	fmt.Scanln(&response)
	return strings.HasPrefix(strings.ToLower(response), "y")
}
