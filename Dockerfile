# Use Go 1.22 as 1.23 is not yet available
FROM docker.io/library/golang:1.23.2-alpine

# Add build dependencies and debugging tools
RUN apk add --no-cache \
    gcc \
    musl-dev \
    git \
    bash \
    curl

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Print Go version for debugging
RUN go version

# Modify go.mod to use Go 1.23.2
RUN go mod edit -go=1.23.2

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -o main ./cmd/api

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"]
