# aws-cf-rtl
AWS Cloudfront Real-Time Logging

## Why?
[AWS Cloudfront](https://aws.amazon.com/cloudfront/) easily stashes JSON formatted weblogs to S3 buckets ([standard logging](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/AccessLogs.html)). However, processing raw JSON files from an S3 bucket is tedious. One can manually update [Athena](https://docs.aws.amazon.com/athena/latest/ug/cloudfront-logs.html) tables to query the data or write custom fetcher/parsers to query the JSON files. To solve the issue, Cloudfront offers a [Real-time logging](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/real-time-logs.html) feature that allows you to stream logs to [Kinesis](https://aws.amazon.com/kinesis/) in real-time. The downside, however, is that setting up the stack from scratch is a bit of a hassle. This repo aims to provide a nearly one-shot [Cloudformation](https://aws.amazon.com/cloudformation/) template to set up and run the Real-Time Logging feature.

## What's in the Box?
This repo provides an AWS Cloudformation template to stand up a basic Cloudfront Real-Time Logging (RTL) service. Items included in the repo:
* Cloudformation template.
* Cloudfront Real-Time Logging (RTL) configuration.
* [AWS Glue](https://aws.amazon.com/glue/) database, table, and crawler.
* Kinesis stream and [Firehose](https://aws.amazon.com/kinesis/data-firehose/) delivery stream (with output conversion to [ORC](https://orc.apache.org)).
* [AWS Lambda](https://aws.amazon.com/lambda/) function to process raw Cloudfront logs into a Glue table-compatible JSON format.
* Basic IAM roles and policies. Note: **THESE ROLES AND POLICES ARE NOT PRODUCTION-READY**.
* S3 bucket for storing raw and processed ORC formatted logs.
* Helper CLI tools: 
  * Process user-agent strings into browser and device type.
  * Process IP addresses into GeoIP data.
  * Raw log re-drive to Kinetisis stream.

## Assumtions: things you should already know or have.
* You have an AWS account with Cloudfront distributions already deployed.
* [Go](https://go.dev) >= 1.17 installed and configured.
* Some level of experience editing AWS Cloudformation templates.
* Be aware: any changes to the data fields selected for real-time logging must be reflected in the [Lambda](./lambda/cf-rtl-kinesis/main.go) function code and the Glue table schema defined in the template.

## Getting Started
* Edit [aws-cloudformation/template.yaml](./aws-cloudformation/template.yaml) to suit your needs. At a minimum, you should edit/verify the `Parameters` section.
* Review IAM Policies and Roles and edit to suit your needs.
* Review and edit the [Makefile](./Makefile), adjusting parameters for your environment. In general, all Cloudfront activites take place in AWS region `us-east-1`.
* Run `make deploy` to build the Lambda function deploy the Cloudformation template.
* Add Cloudfront distributions to the Real-Time Logging service (see [Real-time logs](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/real-time-logs.html) for more information).
* Inkoke hits to the Cloudfront distribution(s).
* Wait at least five minutes for the logs to be processed. Check Cloudwatch logs execution results and errors.
* Check S3 bucket for backup and processed files.

## Next steps
Once you have a full configuration deployed and functional, you can run the provided Glue Crawler to process the ORC formatted logs. Next, use Athena or Trino to query the Glue table.

## Useful Queries
This gist [Useful Trino Queries](https://gist.github.com/rmrfslashbin/b13a37be9aba9266943d42050ef6c74d) provides some useful Trino/Athena queries related to the data stored in the ORC files.

## Feedback
Feedback, comments, pull requests, and questions are welcome.
