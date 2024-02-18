# Stage 1: Building the application
FROM golang:1.21 AS builder

# Copying the entire code and setting work directory
COPY ${PWD} /app
WORKDIR /app

# Download all modules
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Toggle CGO based on your app requirement. CGO_ENABLED=1 for enabling CGO
RUN CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o /app/logAssignment *.go

# Stage 2: Pushing binary to a new image
FROM debian:stable-slim

# Maintainer
LABEL MAINTAINER Teja Surisetty <teja@surisetty.dev>

# Following commands are for installing CA ce   rts (for proper functioning of HTTPS and other TLS)
RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

# Copying the binary from builder stage and setting work directory
COPY --from=builder /app/logAssignment /home/appuser/app/logAssignment
WORKDIR /home/appuser/app

# Exposing the port
EXPOSE 4500

# Entrypoint
CMD ["./logAssignment"]