# Stage 1: compile code.
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download Go depedencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and compile statically.
COPY server/ ./server/
RUN CGO_ENABLED=0 go build --ldflags="-w -s" -o chat-server ./server

# Stage 2: create final image.
FROM alpine:latest

# Create non-root user and group.
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# Copy executable from builder stage.
COPY --from=builder /app/chat-server ./

# Change owner and group of cwd to non-root.
RUN chown -R app:app /app

# Run the image as non-root user.
USER app

# Server listens on port 5001 (defined in code).
EXPOSE 5001

ENV GIN_MODE=release

# Entrypoint that starts the server.
CMD [ "./chat-server" ]
