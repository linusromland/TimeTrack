# --- REACT BUILD STAGE ---
FROM node:20-alpine AS web-builder

WORKDIR /app

# Copy package files for web app
COPY ./apps/web/package.json ./apps/web/package-lock.json* ./apps/web/
RUN cd ./apps/web && npm install

# Copy web app source and build
COPY ./apps/web ./apps/web
RUN cd ./apps/web && npm run build

# --- GO BUILD STAGE ---
FROM golang:1.25.1 AS api-builder

WORKDIR /app

COPY ./apps/api/go.mod ./apps/api/go.sum ./apps/api/
COPY ../../shared ./shared

RUN cd ./apps/api/ && go mod download

COPY ./apps/api/ ./apps/api/

ARG VERSION=dev

RUN cd ./apps/api/ && go build -ldflags "-X 'TimeTrack-api/src/config.Version=${VERSION}'" -o server ./src/main.go

# --- FINAL STAGE ---
FROM ubuntu:24.04

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy Go API server
COPY --from=api-builder /app/apps/api/server ./server

# Copy React static files
COPY --from=web-builder /app/apps/web/dist ./static

ENV GIN_MODE=release
ENV PORT=8080

EXPOSE 8080

CMD ["./server"]
