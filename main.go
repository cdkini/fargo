package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	query, source, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}
	paths, err := getFilesToSearch(source)
	if err != nil {
		log.Fatal(err)
	}
	searchFileList(paths, query)
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

func getFilesToSearch(source string) ([]string, error) {
	f, err := os.Stat(source)
	if err != nil {
		return nil, err
	}

	if f.IsDir() {
		return getFilePathsFromDir(source)
	} else {
		return []string{source}, nil
	}
}

func getFilePathsFromDir(dir string) ([]string, error) {
	paths := make([]string, 0)
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			paths = append(paths, path)
			return nil
		})
	return paths, err
}

func searchFileList(paths []string, query string) {
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
		match := r.MatchString(text)
		if match {
			fmt.Println(path, text)
		}
	}
}
