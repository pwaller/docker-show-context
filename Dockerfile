FROM golang:1.11

WORKDIR /app

COPY . .

RUN GO111MODULE=on go install -v

ENTRYPOINT "docker-show-context"