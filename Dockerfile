FROM golang:1.23-alpine AS build

# Set destination for COPY
WORKDIR /src

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY web ./web
COPY internal ./internal
COPY cmd/app ./cmd/app
COPY embed.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd/app

# Use small image for the final stage
FROM alpine:3.21

# Copy the binary from the build stage
COPY --from=build /server .

# Map to this port to access the server
EXPOSE 8080

# Run
CMD ["/server"]