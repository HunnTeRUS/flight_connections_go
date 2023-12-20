install-deps:
	which go || (curl -O https://golang.org/dl/go1.17.2.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.17.2.linux-amd64.tar.gz)

	which aws || curl "https://d1vvhvl2y92vvt.cloudfront.net/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && unzip awscliv2.zip && sudo ./aws/install

	which sam || (curl -O https://github.com/aws/aws-sam-cli/releases/latest/download/aws-sam-cli-linux-x86_64.zip && unzip aws-sam-cli-linux-x86_64.zip && sudo ./install)

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o ryanairflights .

package:
	zip ryanairflights.zip ryanairflights

upload:
	aws s3 cp ryanairflights.zip s3://ryanair-golang-lambda-test-otavio/

deploy:
	sam deploy --template-file template.yml --stack-name ryanair-flight-connect-stack --force-upload --capabilities CAPABILITY_IAM

start-lambda-locally:
	sam local start-lambda

invoke-local:
	payload_base64=$(echo -n '{"departure": "DUB", "departureDateTime": "2023-11-17T18:00:00Z", "arrival": "WRO", "arrivalDateTime": "2023-11-17T21:50:00Z"}' | base64); \
	aws lambda invoke --function-name ryanair-flight-connections-develop --endpoint-url http://127.0.0.1:3001 --payload "$$payload_base64" output.txt
