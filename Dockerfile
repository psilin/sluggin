FROM golang:1.16-alpine as builder

ENV GIN_MODE=release
ENV PORT=8080

WORKDIR /opt/gin

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY . . 
RUN go build -o ./cmd/server/server ./cmd/server/main.go

EXPOSE $PORT

ENTRYPOINT ["./cmd/server/server", "-c", "./config"]
