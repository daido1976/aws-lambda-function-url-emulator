```sh
$ docker compose up -d
$ go run ./main.go
$ curl "http://localhost:8080/hoge/fuga?foo=bar"
# {"statusCode":200,"headers":{"Content-Type":"application/json"},"body":"{\"message\":\"Hello from Lambda!\",\"input\":{\"version\":\"2.0\",\"rawPath\":\"/hoge/fuga\",\"rawQueryString\":\"foo=bar\",\"headers\":{\"Accept\":\"*/*\",\"User-Agent\":\"curl/8.7.1\"},\"queryStringParameters\":{\"foo\":\"bar\"},\"requestContext\":{\"http\":{\"method\":\"GET\",\"path\":\"/hoge/fuga\"}},\"isBase64Encoded\":false}}"}%
```
