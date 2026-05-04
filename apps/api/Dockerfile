# --- WEB BUILD STAGE ---
FROM node:22-alpine AS web-builder

WORKDIR /app/apps/web

COPY ./apps/web/package.json ./apps/web/package-lock.json ./
RUN npm ci --ignore-scripts

COPY ./apps/web/ ./
RUN npm run build

# --- API BUILD STAGE ---
FROM golang:1.25.1 AS builder

WORKDIR /app

COPY ./apps/api/go.mod ./apps/api/go.sum ./apps/api/
COPY ../../shared ./shared

RUN cd ./apps/api/ && go mod download

COPY ./apps/api/ ./apps/api/

# Copy the built web UI assets so they are embedded into the Go binary
COPY --from=web-builder /app/apps/web/dist ./apps/api/src/web/dist

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
