# Use a base image with Go
FROM golang:1.22.2-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o myapp

# Expose the port your app runs on
EXPOSE 8080

# Run the app
CMD ["./myapp"]