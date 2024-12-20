# AWS Lambda Function URL Emulator

This repository provides a lightweight Docker image to emulate **AWS Lambda Function URLs** locally. It uses **AWS Lambda Runtime Interface Emulator (RIE)** to replicate the behavior of Function URLs for local development and testing.

---

## Features

- Emulates **AWS Lambda Function URLs** locally.
- Automatically forwards HTTP requests to a locally running Lambda function using the AWS Lambda Runtime Interface Emulator.
- Supports `APIGatewayProxyEventV2` for HTTP API requests.
- Handles `isBase64Encoded` for binary data.

---

## Getting Started

### Pull the Docker Image

```bash
docker pull daido1976/aws-lambda-function-url-emulator:latest
```

### Run the Emulator

```bash
docker run --rm -p 8080:8080 \
  -e RIE_ENDPOINT=http://localhost:9000/2015-03-31/functions/function/invocations \
  daido1976/aws-lambda-function-url-emulator:latest
```

The emulator will be available at `http://localhost:8080`.

---

## Usage

1. **Start the Emulator**

   Run the Docker container as described above.

2. **Start a Local Lambda Function with RIE**

   Use **AWS Lambda Runtime Interface Emulator (RIE)** to run your Lambda function locally. Example:

   ```bash
   docker run --rm -v $(pwd):/var/task \
     -p 9000:8080 \
     amazon/aws-lambda-nodejs:20 \
     node index.js
   ```

   Replace `node index.js` with your Lambda handler code.

3. **Send a Request to the Emulator**

   Use `curl`, Postman, or any HTTP client to test your Lambda function through the emulator.

   Example:

   ```bash
   curl http://localhost:8080/?key=example-file.txt
   ```

---

## Environment Variables

| Variable       | Description                                         | Default Value                                                     |
| -------------- | --------------------------------------------------- | ----------------------------------------------------------------- |
| `RIE_ENDPOINT` | URL for the Lambda Runtime Interface Emulator (RIE) | `http://localhost:9000/2015-03-31/functions/function/invocations` |

---

## Build the Docker Image Locally

If you want to build the Docker image yourself:

```bash
docker build -t aws-lambda-function-url-emulator .
```

---

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
