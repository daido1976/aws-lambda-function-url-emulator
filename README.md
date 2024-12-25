# AWS Lambda Function URL Emulator

This repository provides a lightweight Docker image to emulate **AWS Lambda Function URLs** locally. It uses **AWS Lambda Runtime Interface Emulator (RIE)** to replicate the behavior of Function URLs for local development and testing.

## Features

- Emulates **AWS Lambda Function URLs** locally.
- Automatically forwards HTTP requests to a locally running Lambda function using the AWS Lambda Runtime Interface Emulator.
- Supports `APIGatewayProxyEventV2` for HTTP API requests.
- Handles `isBase64Encoded` for binary data.

## Getting Started

### Pull the Docker Image

```bash
docker pull daido1976/aws-lambda-function-url-emulator:latest
```

### Run the Emulator

```bash
docker run --rm -p 8080:8080 \
  -e RIE_ENDPOINT=http://localhost:8080/2015-03-31/functions/function/invocations \
  daido1976/aws-lambda-function-url-emulator:latest
```

The emulator will be available at `http://localhost:8080`.

### Use Docker Compose to work with RIE

For a practical example of integrating this emulator with RIE using Docker Compose, refer to the [example](./example/) directory.

## Environment Variables

| Variable       | Description                                         | Default Value                                                     |
| -------------- | --------------------------------------------------- | ----------------------------------------------------------------- |
| `RIE_ENDPOINT` | URL for the Lambda Runtime Interface Emulator (RIE) | `http://localhost:8080/2015-03-31/functions/function/invocations` |

## Build the Docker Image Locally

If you want to build the Docker image yourself:

```bash
docker build -t aws-lambda-function-url-emulator .
```

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
