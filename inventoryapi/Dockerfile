# Start from the official Go image
FROM golang:1.23
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go.mod and go.sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Expose port 8080 to the outside world
EXPOSE 8080
# Run the Go app
CMD ["go", "run", "./cmd/api/main.go"]
