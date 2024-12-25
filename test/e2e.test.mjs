import { expect, test } from "vitest";

test("E2E Test", async () => {
  const response = await fetch(
    "http://localhost:8080/foo/bar?testkey=testvalue"
  );
  const responseData = await response.json();

  // Expected JSON (dynamic fields replaced with expectations)
  const expected = {
    message: "Hello from Lambda!",
    event: {
      version: "2.0",
      routeKey: "",
      rawPath: "/foo/bar",
      rawQueryString: "testkey=testvalue",
      headers: expect.objectContaining({
        Accept: "*/*",
        "User-Agent": expect.any(String), // User-Agent の形式を検証
      }),
      queryStringParameters: {
        testkey: "testvalue",
      },
      requestContext: {
        routeKey: "",
        accountId: "",
        stage: "",
        requestId: expect.any(String),
        apiId: "",
        domainName: "",
        domainPrefix: "",
        time: expect.any(String),
        timeEpoch: expect.any(Number),
        http: expect.objectContaining({
          method: "GET",
          path: "/foo/bar",
          protocol: "HTTP/1.1",
          sourceIp: expect.any(String),
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

  // Validate the response against the expected structure
  expect(responseData).toMatchObject(expected);
});
