Example of using the `aws-lambda-function-url-emulator` and RIE with Docker Compose.

```sh
$ docker compose up -d
$ curl 'http://localhost:8080/foo/bar?testkey=testvalue'
# {"message":"Hello from Lambda!","event":{"version":"2.0","routeKey":"","rawPath":"/foo/bar","rawQueryString":"testkey=testvalue","headers":{"Accept":"*/*","User-Agent":"curl/8.7.1"},"queryStringParameters":{"testkey":"testvalue"},"requestContext":{"routeKey":"","accountId":"","stage":"","requestId":"","apiId":"","domainName":"","domainPrefix":"","time":"","timeEpoch":0,"http":{"method":"GET","path":"/foo/bar","protocol":"HTTP/1.1","sourceIp":"172.18.0.1:65364","userAgent":"curl/8.7.1"},"authentication":{"clientCert":{"clientCertPem":"","issuerDN":"","serialNumber":"","subjectDN":"","validity":{"notAfter":"","notBefore":""}}}},"isBase64Encoded":false},"context":{"callbackWaitsForEmptyEventLoop":true,"functionVersion":"$LATEST","functionName":"test_function","memoryLimitInMB":"3008","logGroupName":"/aws/lambda/Functions","logStreamName":"$LATEST","invokedFunctionArn":"arn:aws:lambda:us-east-1:012345678912:function:test_function","awsRequestId":"d4908d95-5884-435a-988d-3a67835d9210"}}
```
