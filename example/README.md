Example of using the `aws-lambda-function-url-emulator` and RIE with Docker Compose.

```sh
$ docker compose up -d
$ curl 'http://localhost:8080/foo/bar?testkey=testvalue'
# {"message":"Hello from Lambda!","event":{"version":"2.0","routeKey":"$default","rawPath":"/foo/bar","rawQueryString":"testkey=testvalue","headers":{"Accept":"*/*","User-Agent":"curl/8.7.1"},"queryStringParameters":{"testkey":"testvalue"},"requestContext":{"routeKey":"$default","accountId":"","stage":"$default","requestId":"","apiId":"","domainName":"localhost:8080","domainPrefix":"","time":"26/Dec/2024:01:27:25 +0000","timeEpoch":1735176445,"http":{"method":"GET","path":"/foo/bar","protocol":"HTTP/1.1","sourceIp":"xxx","userAgent":"curl/8.7.1"},"authentication":{"clientCert":{"clientCertPem":"","issuerDN":"","serialNumber":"","subjectDN":"","validity":{"notAfter":"","notBefore":""}}}},"isBase64Encoded":false},"context":{"callbackWaitsForEmptyEventLoop":true,"functionVersion":"$LATEST","functionName":"test_function","memoryLimitInMB":"3008","logGroupName":"/aws/lambda/Functions","logStreamName":"$LATEST","invokedFunctionArn":"arn:aws:lambda:us-east-1:012345678912:function:test_function","awsRequestId":"142e22ca-ae7d-4ada-bb5c-4e3a646c70f1"}}
```
