FROM golang:1.20.1-alpine AS builder
WORKDIR /src/
COPY . .
RUN go get -d -v ./...
RUN go build -o liquidator

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /src/
COPY --from=builder /src/liquidator .

EXPOSE 9000

ENTRYPOINT ["./liquidator"]
