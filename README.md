# Ryanair Flights Lambda

This AWS Lambda function receives a payload to find direct or indirect flights between specified departure and arrival airports within a specified time range.

## Payload Example

```json
{
  "departure": "DUB",
  "departureDateTime": "2023-11-17T18:00:00Z",
  "arrival": "WRO",
  "arrivalDateTime": "2023-11-17T21:50:00Z"
}
```

## Execution

### Execution locally

To run the Lambda function locally using the AWS SAM CLI, follow these steps:

1. Install dependencies (Go, AWS CLI, SAM CLI):

    ```bash
    make install-deps
    ```
   This installs Go, AWS CLI, and SAM CLI if not already installed.


2. Go into template.yml and change the line 29 to:

    ```bash
    CodeUri: .
    ```
   This will force the SAM to use the local code


3. Then, run:

    ```bash
    make start-lambda-locally
    ```
   This start the lambda in your local machine


4. For calling some requestsm do:

    ```bash
    make invoke-local
    ```
   This start the lambda in your local machine


### Execution with SAM

To run the Lambda function using the AWS SAM CLI, follow these steps:

1. Install dependencies (Go, AWS CLI, SAM CLI):

    ```bash
    make install-deps
    ```
   This installs Go, AWS CLI, and SAM CLI if not already installed.


2. Build the Go application:

    ```bash
    make build
    ```
   This will build the golang project and generate a linux executable


3. Package the application:

    ```bash
    make package
    ```
   This will create a ZIP file named ryanairflights.zip in the project root directory.


4. Upload the packaged application to S3:

    ```bash
    make upload
    ```
   This will upload the package to a public S3 bucket named ryanair-golang-lambda-test-otavio


5. Deploy the application using SAM:

    ```bash
    make deploy
    ```
   This uses the SAM CLI to deploy the function and create the necessary AWS resources.

## Expected Output

The Lambda function returns a response in the following format:

```json
[
  {
    "stops": 0,
    "legs": [
      {
        "departureAirport": "DUB",
        "arrivalAirport": "WRO",
        "departureDateTime": "2023-11-17T18:10",
        "arrivalDateTime": "2023-11-17T21:40"
      }
    ]
  }
]
```

This response includes the number of stops ("stops") and details about the available flights.

Note: Ensure that you have the necessary AWS credentials configured for deployment.

## Cleanup
To remove the deployed resources, run the following AWS CLI command:
```bash
aws cloudformation delete-stack --stack-name ryanair-flight-connect-stack
```