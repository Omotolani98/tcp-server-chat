# Use a lightweight Golang image as the base
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

## Copy the Go module files
#COPY go.mod go.sum ./

## Download the dependencies
#RUN go mod download

# Copy the application source code
COPY . .

# Build the application
RUN go build -o tcp-chat-server .

# Expose the port the server runs on
EXPOSE 8080

# Start the server
CMD ["./tcp-chat-server"]