FROM golang:1.18

RUN go version
WORKDIR /tgbot

COPY ./ ./
RUN go install github.com/githubnemo/CompileDaemon@latest
RUN go mod tidy
RUN go mod download

# ENV ENV_FILE=/tgbot/config/docker.env

ENTRYPOINT CompileDaemon --build="go build -o tgbot ./cmd/tarantool-test-task/main.go" --command=./tgbot