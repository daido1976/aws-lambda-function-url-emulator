// @ts-check
import { expect, test } from "vitest";

const TEST_ENDPOINT = "http://localhost:8080/foo/bar";

test("E2E Test for GET Request", async () => {
  const response = await fetch(`${TEST_ENDPOINT}?testkey=testvalue`);
  const responseData = await response.json();

  // Expected JSON for GET request
  const expected = {
    message: "Hello from Lambda!",
    event: {
      version: "2.0",
      routeKey: "$default",
      rawPath: "/foo/bar",
      rawQueryString: "testkey=testvalue",
      headers: expect.objectContaining({
        Accept: "*/*",
        "User-Agent": expect.any(String),
      }),
      queryStringParameters: {
        testkey: "testvalue",
      },
      requestContext: {
        routeKey: "$default",
        stage: "$default",
        domainName: "localhost:8080",
        time: expect.stringMatching(
          /^\d{2}\/[A-Za-z]{3}\/\d{4}:\d{2}:\d{2}:\d{2} [-+]\d{4}$/
        ),
        timeEpoch: expect.any(Number),
        http: expect.objectContaining({
          method: "GET",
          path: "/foo/bar",
          protocol: "HTTP/1.1",
          userAgent: expect.any(String),
        }),
        authentication: expect.objectContaining({
          clientCert: expect.any(Object),
        }),
      },
      isBase64Encoded: false,
    },
    context: {
      callbackWaitsForEmptyEventLoop: true,
      functionVersion: "$LATEST",
      functionName: "test_function",
      memoryLimitInMB: "3008",
      logGroupName: "/aws/lambda/Functions",
      logStreamName: "$LATEST",
      invokedFunctionArn:
        "arn:aws:lambda:us-east-1:012345678912:function:test_function",
      awsRequestId: expect.any(String),
    },
  };

  // Validate the response against the expected structure for GET
  expect(responseData).toMatchObject(expected);
});

test("E2E Test for POST Request", async () => {
  const response = await fetch(TEST_ENDPOINT, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
    },
    body: JSON.stringify({ jsonkey: "jsonvalue" }),
  });
  const responseData = await response.json();

  // Expected JSON for POST request
  const expected = {
    message: "Hello from Lambda!",
    event: {
      version: "2.0",
      routeKey: "$default",
      rawPath: "/foo/bar",
      rawQueryString: "",
      headers: expect.objectContaining({
        Accept: "application/json",
        "Content-Type": "application/json",
        "Content-Length": expect.any(String),
        "User-Agent": expect.any(String),
      }),
      requestContext: {
        routeKey: "$default",
        stage: "$default",
        domainName: "localhost:8080",
        time: expect.stringMatching(
          /^\d{2}\/[A-Za-z]{3}\/\d{4}:\d{2}:\d{2}:\d{2} [-+]\d{4}$/
        ),
        timeEpoch: expect.any(Number),
        http: expect.objectContaining({
          method: "POST",
          path: "/foo/bar",
          protocol: "HTTP/1.1",
          userAgent: expect.any(String),
        }),
        authentication: expect.objectContaining({
          clientCert: expect.any(Object),
        }),
      },
      body: '{"jsonkey":"jsonvalue"}',
      isBase64Encoded: false,
    },
    context: {
      callbackWaitsForEmptyEventLoop: true,
      functionVersion: "$LATEST",
      functionName: "test_function",
      memoryLimitInMB: "3008",
      logGroupName: "/aws/lambda/Functions",
      logStreamName: "$LATEST",
      invokedFunctionArn:
        "arn:aws:lambda:us-east-1:012345678912:function:test_function",
      awsRequestId: expect.any(String),
    },
  };

  // Validate the response against the expected structure for POST
  expect(responseData).toMatchObject(expected);
});

test("E2E Test for Base64 Encoded Response (Binary)", async () => {
  const response = await fetch(`${TEST_ENDPOINT}?binary=true`);
  const responseData = await response.arrayBuffer();

  // Verify binary data after decoding
  const expectedBinary = Buffer.from("This is a test binary data", "utf-8");
  expect(Buffer.from(responseData)).toEqual(expectedBinary);

  // Verify the response headers
  expect(response.headers.get("Content-Type")).toBe("application/octet-stream");
});
