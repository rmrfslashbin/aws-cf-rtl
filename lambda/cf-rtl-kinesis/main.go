package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"github.com/ua-parser/uap-go/uaparser"
)

var (
	// Global logger
	log *logrus.Logger
)

// Record represents a single log entry
type Record struct {
	Timestamp                int64   `json:"timestamp"`
	ClientIP                 net.IP  `json:"client_ip"`
	Status                   int     `json:"status"`
	Bytes                    int64   `json:"bytes"`
	Method                   string  `json:"method"`
	Protocol                 string  `json:"protocol"`
	Host                     string  `json:"host"`
	URIStem                  string  `json:"uri_stem"`
	EdgeLocation             string  `json:"edge_location"`
	EdgeRequestId            string  `json:"edge_request_id"`
	HostHeader               string  `json:"host_header"`
	TimeTaken                float64 `json:"time_taken"`
	ProtoVersion             string  `json:"proto_version"`
	IPVersion                string  `json:"ip_version"`
	UserAgent                string  `json:"user_agent"`
	Referer                  string  `json:"referer"`
	Cookie                   string  `json:"cookie"`
	URIQuery                 string  `json:"uri_query"`
	EdgeResponseResultType   string  `json:"edge_response_result_type"`
	SSLProtocol              string  `json:"ssl_protocol"`
	SSLCipher                string  `json:"ssl_cipher"`
	EdgeResultType           string  `json:"edge_result_type"`
	ContentType              string  `json:"content_type"`
	ContentLength            int64   `json:"content_length"`
	EdgeDetailedResultType   string  `json:"edge_detailed_result_type"`
	Country                  string  `json:"country"`
	CacheBehaviorPathPattern string  `json:"cache_behavior_path_pattern"`
	UserAgentDeviceFamily    string  `json:"user_agent_device_family"`
	UserAgentDeviceBrand     string  `json:"user_agent_device_brand"`
	UserAgentDeviceModel     string  `json:"user_agent_device_model"`
	UserAgentOSFamily        string  `json:"user_agent_os_family"`
	UserAgentOSMajor         string  `json:"user_agent_os_major"`
	UserAgentOSMinor         string  `json:"user_agent_os_minor"`
	UserAgentOSPatch         string  `json:"user_agent_os_patch"`
	UserAgentOSPatchMinor    string  `json:"user_agent_os_patch_minor"`
	UserAgentFamily          string  `json:"user_agent_family"`
	UserAgentMajor           string  `json:"user_agent_major"`
	UserAgentMinor           string  `json:"user_agent_minor"`
	UserAgentPatch           string  `json:"user_agent_patch"`
}

/* Fields mapped from Cloudfront Real-Time Logs configuration.

00 timestamp (string) (len=14) "1642349408.581",
01 c-ip	(string) (len=14) "123.123.123.123",
02 sc-status (string) (len=3) "200",
03 sc-bytes (string) (len=4) "3536",
04 cs-method (string) (len=3) "GET",
05 cs-protocol (string) (len=5) "https",
06 cs-host (string) (len=27) "www.example.com",
07 cs-uri-stem (string) (len=13) "/news/today/",
08 x-edge-location (string) (len=8) "IAD89-P2",
09 x-edge-request-id (string) (len=56) "fv2T3ZdTRe4x0VV4Ro6YLWhfvD0LvfeKVRtJAXXWaev6SxFOPjhkjM==",
10 x-host-header (string) (len=29) "d986b4ld3rmrlc.cloudfront.net",
11 time-taken (string) (len=5) "0.130",
12 cs-protocol-version (string) (len=8) "HTTP/1.1",
13 c-ip-version (string) (len=4) "IPv4",
14 cs-user-agent (string) (len=83) "Mozilla/5.0%20(compatible;%20SemrushBot/7%7Ebl;%20+http://www.semrush.com/bot.html)",
15 cs-referer (string) (len=1) "-",
16 cs-cookie (string) (len=1) "-",
17 cs-uri-query (string) (len=1) "-",
18 x-edge-response-result-type (string) (len=4) "Miss",
19 ssl-protocol (string) (len=7) "TLSv1.3",
20 ssl-cipher (string) (len=22) "TLS_AES_128_GCM_SHA256",
21 x-edge-result-type (string) (len=4) "Miss",
22 sc-content-type (string) (len=9) "text/html",
23 sc-content-len (string) (len=1) "-",
24 x-edge-detailed-result-type (string) (len=4) "Miss",
25 c-country (string) (len=2) "GB",
26 cache-behavior-path-pattern (string) (len=1) "*"
*/

// handler is the Lambda function handler
func handler(ctx context.Context, kinesisFirehoseEvent events.KinesisFirehoseEvent) (*events.KinesisFirehoseResponse, error) {
	// Struct to hold the response
	output := &events.KinesisFirehoseResponse{}

	// Loop through each record from Kinesis Firehose
	for _, record := range kinesisFirehoseEvent.Records {
		// Split the data by tabs
		parts := strings.Split(string(record.Data), "\t")

		// Marshal the data into a Record struct
		data, err := marshal(&parts)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Error("marshal failed")
			continue
		}

		// Convert to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Error("marshal failed")
			continue
		}

		// Add partition keys to the response
		partitionKeys := make(map[string]string)
		timeMili := time.UnixMilli(data.Timestamp)
		partitionKeys["year"] = fmt.Sprintf("%d", timeMili.Year())
		partitionKeys["month"] = fmt.Sprintf("%d", timeMili.Month())
		partitionKeys["day"] = fmt.Sprintf("%d", timeMili.Day())

		// Create the response
		output.Records = append(output.Records, events.KinesisFirehoseResponseRecord{
			RecordID: record.RecordID,
			Result:   events.KinesisFirehoseTransformedStateOk,
			Data:     jsonData,
			Metadata: events.KinesisFirehoseResponseRecordMetadata{
				PartitionKeys: partitionKeys,
			},
		})
	}
	// Return the response to Kinesis Firehose
	return output, nil
}

// init the logger and other things as needed
func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
}

// main is the entry point
func main() {
	// Catch errors
	var err error
	defer func() {
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("main crashed")
		}
	}()

	// Run the lambda function
	lambda.Start(handler)
}

// marshal the log fields into a Record
func marshal(parts *[]string) (*Record, error) {
	record := &Record{}

	if tstampfloat, err := strconv.ParseFloat((*parts)[0], 64); err == nil {
		record.Timestamp = int64(tstampfloat * 1000)
	} else {
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"timestamp": (*parts)[0],
		}).Error("timestamp failed")
	}

	record.ClientIP = net.ParseIP((*parts)[1])
	record.Status, _ = strconv.Atoi((*parts)[2])
	record.Bytes, _ = strconv.ParseInt((*parts)[3], 10, 64)
	record.Method = (*parts)[4]
	record.Protocol = (*parts)[5]
	record.Host = (*parts)[6]
	record.URIStem = (*parts)[7]
	record.EdgeLocation = (*parts)[8]
	record.EdgeRequestId = (*parts)[9]
	record.HostHeader = (*parts)[10]
	record.TimeTaken, _ = strconv.ParseFloat((*parts)[11], 64)
	record.ProtoVersion = (*parts)[12]
	record.IPVersion = (*parts)[13]
	record.UserAgent = (*parts)[14]
	record.Referer = (*parts)[15]
	record.Cookie = (*parts)[16]
	record.URIQuery = (*parts)[17]
	record.EdgeResponseResultType = (*parts)[18]
	record.SSLProtocol = (*parts)[19]
	record.SSLCipher = (*parts)[20]
	record.EdgeResultType = (*parts)[21]
	record.ContentType = (*parts)[22]
	record.ContentLength, _ = strconv.ParseInt((*parts)[23], 10, 64)
	record.EdgeDetailedResultType = (*parts)[24]
	record.Country = (*parts)[25]
	record.CacheBehaviorPathPattern = (*parts)[26]

	// Add user agent parsing data
	parser := uaparser.NewFromSaved()
	client := parser.Parse(record.UserAgent)
	record.UserAgentDeviceFamily = client.Device.Family
	record.UserAgentDeviceBrand = client.Device.Brand
	record.UserAgentDeviceModel = client.Device.Model
	record.UserAgentOSFamily = client.Os.Family
	record.UserAgentOSMajor = client.Os.Major
	record.UserAgentOSMinor = client.Os.Minor
	record.UserAgentOSPatch = client.Os.Patch
	record.UserAgentOSPatchMinor = client.Os.PatchMinor
	record.UserAgentFamily = client.UserAgent.Family
	record.UserAgentMajor = client.UserAgent.Major
	record.UserAgentMinor = client.UserAgent.Minor
	record.UserAgentPatch = client.UserAgent.Patch

	return record, nil
}
