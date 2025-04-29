# Stage 1: Build the Go application
FROM golang:1.24 AS build
WORKDIR /app
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download
COPY ./src/main.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# Stage 2: Create a minimal image for running the application
FROM alpine:latest AS runtime
WORKDIR /app

# Create a non-root user with a specific UID
RUN addgroup -S appgroup && adduser -S -u 1001 appuser -G appgroup

# Copy the built application
COPY --from=build /app/app .

# Change ownership of the application to the non-root user
RUN chown appuser:appgroup /app/app

# Switch to the non-root user
USER appuser

CMD ["./app"]