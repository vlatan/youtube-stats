FROM golang:1.23-alpine AS build

# Set destination for COPY
WORKDIR /src

# Copy the source code
COPY go.mod go.sum embed.go ./
COPY web ./web
COPY internal ./internal
COPY cmd/app ./cmd/app

# Download Go modules
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /binary ./cmd/app

# Use small image for the final stage
FROM alpine:3.21

# Copy the binary from the build stage
COPY --from=build /binary .

# Run
CMD ["/binary"]