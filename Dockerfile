# Use an official Go runtime as a parent image
FROM golang:1.23-alpine

# Add Maintainer Info
LABEL maintainer="MintTnim19 <mint.tnim19@gmail.com>"

# Environment variables which CompileDaemon requires to run
ENV GIN_MODE=release

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download && go mod verify

# Copy the source code into the container 
COPY . .

# Build the Go app
RUN cp .env.production .env && rm .env.* && go build -o /usr/local/bin/app ./cmd/server/main.go

# Command to run the executable
CMD ["app"]
