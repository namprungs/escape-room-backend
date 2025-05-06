# Stage 1: Build
FROM golang:1.23 as builder

WORKDIR /app

# Copy go files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build binary (main.go is in cmd/)
RUN go build -o server ./cmd/main.go

# Stage 2: Run
FROM gcr.io/distroless/base-debian11

WORKDIR /app
COPY --from=builder /app/server /app/server

# Expose port (matching APP_PORT)
EXPOSE 5000

# Run binary
CMD ["/app/server"]
