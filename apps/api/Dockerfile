# --- BUILD STAGE ---
FROM golang:1.25.1 AS builder

WORKDIR /app

COPY ./apps/api/go.mod ./apps/api/go.sum ./apps/api/
COPY ../../shared ./shared

RUN cd ./apps/api/ && go mod download

COPY ./apps/api/ ./apps/api/

ARG VERSION=dev

RUN cd ./apps/api/ && go build -ldflags "-X 'TimeTrack-api/src/config.Version=${VERSION}'" -o server ./src/main.go

FROM ubuntu:24.04

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

COPY --from=builder /app/apps/api/server ./server

ENV GIN_MODE=release
ENV PORT=8080

EXPOSE 8080

CMD ["./server"]
