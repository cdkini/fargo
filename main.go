package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	args := ParseArgs()
	results := RunRipgrep(args)
	filesToChange := FilterRelevantResults(results)
	fmt.Println(filesToChange)
}

func ParseArgs() []string {
	args := os.Args[1:]
	args = append(args, "-n") // Line numbers are essential for proper SearchResult parsing
	return args
}

type SearchResult struct {
	path    string
	line    int
	text    string
	indices []int
}

func RunRipgrep(args []string) []SearchResult {
	rg := exec.Command("rg", args...)
	out, err := rg.Output()
	if err != nil {
		log.Fatal(err)
	}

	results := make([]SearchResult, 0)
	lines := strings.Split(string(out), "\n")

	// <FILE_PATH>:<LINE_NUMBER>:<MATCHING_TEXT>
	rgRegex := regexp.MustCompile(`(.+):(\d+):(.+)`)
	// Whatever pattern provided by user as first arg to rg
	lineRegex := regexp.MustCompile(args[0])

	for _, line := range lines {
		result, err := parseResultFromString(line, lineRegex, rgRegex)
		if err != nil {
			continue // Failed parsing shouldn't cause a sys exit
		}
		results = append(results, result)
	}

	return results
}

func parseResultFromString(str string, lineRegex *regexp.Regexp, rgRegex *regexp.Regexp) (SearchResult, error) {
	matches := rgRegex.FindStringSubmatch(str)
	if len(matches) != 4 {
		err := errors.New("Something went wrong when parsing text")
		return SearchResult{}, err
	}

	path := matches[1]
	line, err := strconv.Atoi(matches[2])
	text := matches[3]
    indices := lineRegex.FindStringIndex(text) // TODO(cdkini): Ensure this captures ALL matches in a line
	if err != nil {
		log.Fatal(err)
	}

	return SearchResult{path, line, text, indices}, nil
}

func FilterRelevantResults(results []SearchResult) []SearchResult {
	relevantResults := make([]SearchResult, 0)
	for _, result := range results {
		if promptUserConfirmation(result) {
			relevantResults = append(relevantResults, result)
		}
	}
	return relevantResults
}

// TODO(cdkini): This needs refinement
func promptUserConfirmation(result SearchResult) bool {
	yellow := color.New(color.FgYellow).Add(color.Underline)
	yellow.Printf("\n%s L%d\n\n", result.path, result.line)

    before := result.text[:result.indices[0]]
    match := result.text[result.indices[0]:result.indices[1]]
    after := result.text[result.indices[1]:] 

    red := color.New(color.FgRed).Add(color.Bold)
    fmt.Print(before)
    red.Print(match)
    fmt.Println(after)

	cyan := color.New(color.FgCyan).Add(color.Bold)
	cyan.Print("\nReplace [y/n]: ")

	var response string
	fmt.Scanln(&response)

	return strings.HasPrefix(strings.ToLower(response), "y")
}
