# Stage 1: Build the Go binary
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# The output binary will be named 'data-analyst-adk' based on the module name.
RUN go build -o /app/data-analyst-adk .

# Stage 2: Create the final lightweight image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/data-analyst-adk .

# Expose the port the application runs on
# The ADK web launcher defaults to 8080
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["./data-analyst-adk"]
CMD ["web", "webui", "-api_server_address", "https://data-analyst-agent-280946129258.asia-southeast1.run.app/api", "api", "-webui_address", "https://data-analyst-agent-280946129258.asia-southeast1.run.app" ]
