# --- BUILD STAGE ---
FROM golang:1.25.1 AS builder

WORKDIR /app

COPY ./apps/api/go.mod ./apps/api/go.sum ./apps/api/
COPY ../../shared ./shared

RUN cd ./apps/api/ && go mod download

COPY ./apps/api/ ./apps/api/

ARG VERSION=dev

RUN cd ./apps/api/ && go build -ldflags "-X 'TimeTrack-api/src/config.AppConfig.APIVersion=${VERSION}'" -o server ./src/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/apps/api/server .

EXPOSE 8080

CMD ["./server"]
