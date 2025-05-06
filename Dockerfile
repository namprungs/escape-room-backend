# ---------- Stage 1: Build ----------
  FROM golang:1.23 as builder

  WORKDIR /app
  
  # Copy go files and download dependencies
  COPY go.mod go.sum ./
  RUN go mod download
  
  # Copy the rest of the source
  COPY . .
  
  # Build binary (main.go is in cmd/)
  RUN go build -o server ./cmd/main.go
  
  # ---------- Stage 2: Run ----------
  FROM alpine:latest
  
  WORKDIR /app
  
  # Copy binary and start script
  COPY --from=builder /app/server /app/server
  COPY start.sh ./start.sh
  
  # Make sure start.sh is executable
  RUN chmod +x start.sh
  
  # Set ENV if needed
  ENV APP_PORT=5000
  
  EXPOSE 5000
  
  # Run using shell
  CMD ["sh", "start.sh"]
  