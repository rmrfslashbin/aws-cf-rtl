package main

// This example shows how to lookup up GeoIP data extracted from logs.

import (
	"bufio"
	"log"
	"net"
	"os"
	"path"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/rmrfslashbin/aws-cf-rtl/pkg/geoip"
)

const (
	datafile = "" // path to txt data file; one line per IP
)

func main() {
	// Sanatinize the file path
	fqpn := path.Clean(datafile)

	// Open the file
	datafile, err := os.Open(fqpn)
	if err != nil {
		panic(err)
	}
	defer datafile.Close()

	// User a scanner to read the file
	scanner := bufio.NewScanner(datafile)
	for scanner.Scan() {
		// Get a line and trim whitespace
		line := strings.TrimSpace(scanner.Text())

		// Parse the line into an IP address
		ip := net.ParseIP(line)

		// Get the geoip data
		data, err := geoip.Lookup(ip)
		if err != nil {
			log.Fatal(err)
		}

		// Display the data
		spew.Dump(data)
	}

	// Check for errors
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
