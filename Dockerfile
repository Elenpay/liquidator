FROM golang:1.25-alpine AS builder
    WORKDIR /src/
    COPY . .
    RUN go get -d -v ./...
    RUN go build -o liquidator

FROM alpine:latest
    RUN adduser -D ep
    RUN apk --no-cache add ca-certificates
    USER ep
    WORKDIR /home/ep
    COPY --from=builder /src/liquidator /usr/bin/liquidator
    EXPOSE 9000
    ENTRYPOINT ["liquidator"]