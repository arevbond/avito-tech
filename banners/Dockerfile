FROM golang:latest AS builder

WORKDIR /app
COPY . .

RUN go get -d -v ./...
RUN go build -o avito-tech ./cmd/avito-tech/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app .

CMD ["./avito-tech", "-config=./cmd/avito-tech/config/config.yaml"]
