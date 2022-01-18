package main

// This CLI tool is used to re-drive backup logs to a Kinesis stream.

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

const (
	//testFile   = "backup/cf-rtl-logs-delivery-stream-1-2022-01-18-00-18-55-1adc337f-84b7-4a04-9620-96890081fec4"
	aws_region = "us-east-1"         // e.g. us-east-1
	rootdir    = "backup"            // Location of backup files
	streamName = "my_kinesis_stream" // Name of Kinesis stream
)

var (
	count int
	kc    *kinesis.Client
)

func init() {
	// The counter is used to display the number of files processed.
	count = 0

	// Create a Kinesis client
	c, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = aws_region
		return nil
	})
	if err != nil {
		panic(err)
	}
	kc = kinesis.NewFromConfig(c)
}

func main() {
	// Process a single file
	//parseFile(testFile)

	// Process all files in a directory
	parseDir(rootdir)
	fmt.Printf("\n%d\n", count)
}

func parseDir(dirname string) {
	// Sanatinize the directory path
	dirfqpn := path.Clean(dirname)

	// Read the directory
	files, err := ioutil.ReadDir(dirfqpn)
	if err != nil {
		panic(err)
	}

	// Loop through the files
	for _, file := range files {
		// Create a clean fully qualified path name for the file
		fqpn := path.Clean(path.Join(dirfqpn, file.Name()))

		// Process the file
		parseFile(fqpn)
	}
}

func parseFile(filename string) {
	// Sanaitize the file path
	fqpn := path.Clean(filename)
	fmt.Println(fqpn)

	// Open the file
	datafile, err := os.Open(fqpn)
	if err != nil {
		panic(err)
	}
	defer datafile.Close()

	// User a scanner to read the file
	scanner := bufio.NewScanner(datafile)
	for scanner.Scan() {
		// Get a line and split by 'tab'
		line := scanner.Text()
		parts := strings.Split(line, "\t")

		// Each log should have 27 fields.
		// Adjust to suit configuration of your log.
		if len(parts) != 27 {
			panic(fmt.Sprintf("Wrong field count: %d", len(parts)))
		}

		// Send the log to Kinesis
		if opt, err := kc.PutRecord(context.TODO(), &kinesis.PutRecordInput{
			Data:         []byte(line),
			StreamName:   aws.String(streamName),
			PartitionKey: aws.String("partitionKey"),
		}); err != nil {
			fmt.Println(line)
			panic(err)
		} else {
			fmt.Printf("  %s :: %s", *opt.SequenceNumber, *opt.ShardId)
		}

		count++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
