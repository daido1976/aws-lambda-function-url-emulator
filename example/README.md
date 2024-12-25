Example of using the `aws-lambda-function-url-emulator` and RIE with Docker Compose.

```sh
$ docker compose up -d
$ curl 'http://localhost:8080/foo/bar?testkey=testvalue'
# {"message":"Hello from Lambda!","input":{"version":"2.0","routeKey":"","rawPath":"/foo/bar","rawQueryString":"testkey=testvalue","headers":{"Accept":"*/*","User-Agent":"curl/8.7.1"},"queryStringParameters":{"testkey":"testvalue"},"requestContext":{"routeKey":"","accountId":"","stage":"","requestId":"","apiId":"","domainName":"","domainPrefix":"","time":"","timeEpoch":0,"http":{"method":"GET","path":"/foo/bar","protocol":"HTTP/1.1","sourceIp":"172.20.0.1:55362","userAgent":"curl/8.7.1"},"authentication":{"clientCert":{"clientCertPem":"","issuerDN":"","serialNumber":"","subjectDN":"","validity":{"notAfter":"","notBefore":""}}}},"isBase64Encoded":false}}%
```
