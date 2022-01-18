package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/rmrfslashbin/aws-cf-rtl/pkg/useragent"
)

const (
	datafile = "" // text file with one user-agent string per line
)

func main() {
	// Sanatize the datafile path.
	fqpn := path.Clean(datafile)
	fmt.Println(fqpn)

	// Open the file.
	datafile, err := os.Open(fqpn)
	if err != nil {
		log.Fatal(err)
	}
	defer datafile.Close()

	// Create a new scanner.
	scanner := bufio.NewScanner(datafile)
	for scanner.Scan() {
		// Remove witespace.
		line := strings.TrimSpace(scanner.Text())

		// Parse the user-agent string.
		records := useragent.Parse(line)

		// Print the results.
		spew.Dump(records)
	}
	// Check for errors.
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
