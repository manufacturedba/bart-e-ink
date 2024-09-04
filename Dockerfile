# Matching version from go.mod
FROM golang:1.23 AS build

# Current working directory
WORKDIR /app

# Copy everything over
COPY . .

# Download all dependencies
RUN go mod download

# Build a static binary
ENV CGO_ENABLED=0

# Build the Go app
RUN go build ./cmd/server

# Expose port 8080 to the outside world
EXPOSE 8080

# Set PORT
ENV PORT=8080

# Run it
CMD ["./server"]