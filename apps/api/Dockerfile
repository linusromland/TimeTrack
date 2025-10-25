# --- BUILD STAGE ---
FROM golang:1.25.1 AS builder

WORKDIR /app

COPY ./apps/api/go.mod ./apps/api/go.sum ./apps/api/
RUN cd ./apps/api/ && go mod download

COPY ./apps/api/ ./apps/api/

COPY ../../shared ./shared

ARG VERSION=dev

RUN go build -ldflags "-X 'TimeTrack-api/src/config.AppConfig.APIVersion=${VERSION}'" -o server ./apps/api/cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
