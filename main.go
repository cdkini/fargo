package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
)

func main() {
	query, source, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}
	files := getFilesToSearch(source)
	searchFileList(files, query)
}

func parseArgs() (string, string, error) {
	argv := flag.Args()
	argc := len(argv)
	if argc == 1 {
		return argv[0], ".", nil
	} else if argc == 2 {
		return argv[0], argv[1], nil
	}
	return "", "", errors.New("A search query/pattern is mandatory")
}

func getFilesToSearch(source string) []string {
	return []string{source}
}

func searchFileList(files []string, query string) {
	for _, file := range files {
		searchFile(file, query)
	}
}

func searchFile(file, query string) {
	fmt.Println(file, query)
}
