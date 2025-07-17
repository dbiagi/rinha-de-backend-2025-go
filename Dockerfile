# Builder
FROM golang:1.24.1-alpine3.21 AS builder

RUN apk update && apk upgrade && \
    apk --update add git make bash build-base

WORKDIR /app

COPY . .

RUN make build

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

EXPOSE 9999

COPY --from=builder /app/bin /app/bin
CMD ["/app/bin/app", "serve"]
