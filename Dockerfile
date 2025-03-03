# Use the official Golang image as the build stage
FROM golang:1.23 AS builder

# Set environment variables to ensure `go install` puts `swag` in a known location
ENV GOPATH=/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOBIN:$PATH

# Set the working directory inside the container
WORKDIR /app

# Install Swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Verify that Swag CLI is installed
RUN swag --version

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -ldflags "-s -w" -o main .

# Use the official Golang image as the base image for the final stage
FROM golang:1.23

# Install necessary runtime dependencies
RUN apt-get update && apt-get install -y gettext-base && rm -rf /var/lib/apt/lists/*

# Install Chromium and its dependencies
RUN apt-get update && apt-get install -y \
  chromium \
  fonts-liberation \
  libappindicator3-1 \
  libasound2 \
  libatk-bridge2.0-0 \
  libatk1.0-0 \
  libcups2 \
  libdbus-1-3 \
  libgdk-pixbuf2.0-0 \
  libnspr4 \
  libnss3 \
  libx11-xcb1 \
  libxcomposite1 \
  libxdamage1 \
  libxrandr2 \
  xdg-utils \
  --no-install-recommends && \
  rm -rf /var/lib/apt/lists/*

# Install Chinese fonts
RUN apt-get update && apt-get install -y \
  fonts-wqy-zenhei \
  --no-install-recommends && \
  rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Create the /storage directory
RUN mkdir -p /storage && chmod -R 777 /storage

# Copy the storage directory
COPY storage /app/storage

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .
COPY config.template.json /app/config.template.json
COPY init-config.sh /app/init-config.sh

# Make the initialization script executable
RUN chmod +x /app/init-config.sh

# Expose the port on which the application will run
EXPOSE 8000

# Command to run the initialization script
CMD ["/app/init-config.sh"]