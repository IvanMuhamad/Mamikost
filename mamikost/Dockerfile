#build stage
FROM golang:alpine AS builder
# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code to the Working Directory inside the container
COPY . .

# Download all dependencies.
# Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o mamikost main.go

# Export necessary ports
EXPOSE 8080

# Command to run when starting the container
CMD ["/app/mamikost"]

