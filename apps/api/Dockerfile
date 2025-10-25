FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG VERSION=dev

RUN go build -ldflags "-X 'TimeTrack-api/src/config.AppConfig.APIVersion=${VERSION}'" -o server ./src/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
