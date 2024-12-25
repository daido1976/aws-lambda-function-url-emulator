# See. https://github.com/GoogleContainerTools/distroless/blob/530158861eebdbbf149f7e7e67bfe45eb433a35c/examples/go/Dockerfile
FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

# Build the Go application with static linking
RUN CGO_ENABLED=0 go build -o main .

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
