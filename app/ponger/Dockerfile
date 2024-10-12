# Define the base image
FROM golang:latest AS builder

# Set working directory for build context
WORKDIR /app

# Copy source code
COPY . .

# Install dependencies
RUN go mod download

# Build the application (change directory if needed)
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a --installsuffix cgo --ldflags="-s -w" -o ponger .

# Define the final image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary
COPY --from=builder /app/ponger .
COPY --from=builder /app/config.yaml .

# Expose the port (adjust if needed)
EXPOSE 8080

# Command to run the application
CMD ["./ponger"]

# Optional: Copy configuration file (adjust path if needed)
# COPY config.yaml /app/config.yaml