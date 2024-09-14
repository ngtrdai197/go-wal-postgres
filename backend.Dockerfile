# Setup Dockerfile for backend by golang version 1.23.6
FROM golang:1.21-alpine3.18 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.18

# Install curl
RUN apk add --no-cache curl

COPY --from=builder /app/main /app/
COPY ./config/config.yml /app/config/config.yml

WORKDIR /app

# Expose port 8088 to the outside world
EXPOSE 8088

ENV GIN_MODE=release

## Command to run the executable
#ENTRYPOINT ["./main public-api"]