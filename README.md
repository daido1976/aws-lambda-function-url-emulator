# AWS Lambda Function URL Emulator

A lightweight Docker image to emulate **AWS Lambda Function URLs** locally. It works with [AWS Lambda Runtime Interface Emulator (RIE)](https://github.com/aws/aws-lambda-runtime-interface-emulator/) to replicate the behavior of Function URLs for local development and testing.

## Features

- Emulates **AWS Lambda Function URLs** locally.
- Automatically forwards HTTP requests to a locally running Lambda function to work with the AWS Lambda Runtime Interface Emulator.
- Supports `APIGatewayProxyEventV2` for HTTP API requests.
- Handles `isBase64Encoded` for binary data.
- Supports enabling CORS via an environment variable, allowing cross-origin requests.

## Getting Started

### Run with Docker (Standalone)

Please ensure that RIE is running in a separate container beforehand.

```bash
docker run --rm -p 8080:8080 \
  -e RIE_ENDPOINT=http://host.docker.internal:9000/2015-03-31/functions/function/invocations \
  daido1976/aws-lambda-function-url-emulator:latest
```

The emulator will be available at `http://localhost:8080` and forwards requests to the RIE endpoint.

### Run with Docker Compose

For a practical example of integrating this emulator with RIE using Docker Compose, refer to the [example](./example/) directory.

### Run as a CLI Tool (Go)

Install the emulator as a CLI tool:

```bash
go install github.com/daido1976/aws-lambda-function-url-emulator@latest
```

Run it with:

```bash
aws-lambda-function-url-emulator
```

The emulator will be available at `http://localhost:8080` and forwards requests to the default RIE endpoint. To use a custom endpoint, set the `RIE_ENDPOINT` environment variable:

```bash
RIE_ENDPOINT=http://custom-host:9000/2015-03-31/functions/function/invocations aws-lambda-function-url-emulator
```

## Environment Variables

| Variable       | Description                                                                                                                                                                                                 | Default Value                                                     |
| -------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| `RIE_ENDPOINT` | URL for the Lambda Runtime Interface Emulator (RIE)                                                                                                                                                         | `http://localhost:9000/2015-03-31/functions/function/invocations` |
| `ENABLE_CORS`  | Set this to `"true"` to enable CORS for all origins, methods, and headers.<br>**Warning:** Enabling this setting may have security implications, Ensure you understand the implications before enabling it. | `"false"`                                                         |

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
