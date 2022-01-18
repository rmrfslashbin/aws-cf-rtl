stack_name = cf-rtl
deploy_bucket = is-us-east-1-deployment
aws_profile = default

.PHONY: build

build:
	GOOS=linux GOARCH=arm64 go build -o bin/cf-rtl-kinesis/bootstrap lambda/cf-rtl-kinesis/main.go
	rm -f bin/cf-rtl-kinesis/bootstrap.zip
	zip -j bin/cf-rtl-kinesis/bootstrap.zip bin/cf-rtl-kinesis/bootstrap
	
deploy: build
	aws --profile $(aws_profile) cloudformation package --template-file aws-cloudformation/template.yaml --s3-bucket $(deploy_bucket) --output-template-file build/out.yaml
	aws --profile $(aws_profile) cloudformation deploy --template-file build/out.yaml --s3-bucket $(deploy_bucket) --stack-name $(stack_name) --capabilities CAPABILITY_NAMED_IAM

