FROM golang:1.11 AS builder
WORKDIR /app
COPY . .
ENV GO111MODULE=on
RUN go build -o docker-show-context .

FROM alpine:latest
COPY --from=builder /app/docker-show-context /usr/local/bin
RUN apk add --no-cache libc6-compat
WORKDIR /data
ENTRYPOINT ["docker-show-context"]
