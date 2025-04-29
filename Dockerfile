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
COPY --from=build /app/app .
CMD ["./app"]