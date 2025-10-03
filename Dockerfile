# syntax=docker/dockerfile:1

# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Dockerfile reference guide at
# https://docs.docker.com/go/dockerfile-reference/

# Want to help us make this template better? Share your feedback here: https://forms.gle/ybq9Krt8jtBL3iCk7

################################################################################
# Create a stage for building the application.
ARG GO_VERSION=1.24.4
FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .

#FROM alpine:latest AS final
RUN go build -o main ./cmd/app/*.go

# Expose the port that the application listens on.
EXPOSE 8080

# What the container should run when it is started.
CMD [ "./main" ]
