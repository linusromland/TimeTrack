# --- BUILD STAGE ---
FROM golang:1.25.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./src/main.go

# --- FINAL STAGE ---
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/app .

# Set Gin to production mode
ENV GIN_MODE=release

EXPOSE 8080

CMD ["./app"]
