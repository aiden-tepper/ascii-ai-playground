FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all the dependencies.
RUN go mod download

# Build the Go app
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o 8ball .
RUN go build -o 8ball .

# Start a new stage from scratch
FROM ubuntu:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/8ball .
RUN chmod +x /root/8ball

# Install gotty
RUN apt-get update && \
    apt-get install -y --no-install-recommends wget ca-certificates && \
    wget https://github.com/sorenisanerd/gotty/releases/download/v1.5.0/gotty_v1.5.0_linux_amd64.tar.gz \
    && tar -xzf gotty_v1.5.0_linux_amd64.tar.gz \
    && mv gotty /usr/local/bin/gotty \
    && rm gotty_v1.5.0_linux_amd64.tar.gz

# Run gotty with your application
CMD ["gotty", "-w", "/root/8ball"]
