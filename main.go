package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	query, source, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	paths, err := GetFilesToSearch(source)
	if err != nil {
		log.Fatal(err)
	}

	SearchFiles(paths, query)
}

func parseArgs() (string, string, error) {
	args := os.Args[1:]
	if len(args) == 1 {
		return args[0], ".", nil
	} else if len(args) == 2 {
		return args[0], args[1], nil
	}
	return "", "", errors.New("Improper arguments; search query/pattern is mandatory")
}

func GetFilesToSearch(source string) ([]string, error) {
	f, err := os.Stat(source)
	if err != nil {
		return nil, err
	}

	if f.IsDir() {
		return getFilePathsFromDir(source)
	}
	return []string{source}, nil
}

func getFilePathsFromDir(dir string) ([]string, error) {
	paths := make([]string, 0)
	err := filepath.Walk(dir,
		func(path string, file os.FileInfo, err error) error {
			if isHiddenDir(file) {
				return filepath.SkipDir
			}
			paths = append(paths, path)
			return nil
		})

	return paths, err
}

func isHiddenDir(file os.FileInfo) bool {
	return file.IsDir() && strings.HasPrefix(file.Name(), ".")
}

func SearchFiles(paths []string, query string) {
	r, err := regexp.Compile(query)
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range paths {
		searchFile(path, r)
	}
}

func searchFile(path string, r *regexp.Regexp) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if match := r.MatchString(text); match {
			fmt.Println(path)
		}
	}
}
