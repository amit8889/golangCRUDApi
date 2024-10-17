FROM golang:1.23-alpine AS build
WORKDIR /app

# Install build dependencies (including gcc for CGO)
RUN apk add --no-cache gcc musl-dev

# Set CGO_ENABLED to 1 so sqlite3 can work
ENV CGO_ENABLED=1

# Copy necessary files and build the Go binary
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o students-api ./cmd/students-api

FROM alpine:3.18
WORKDIR /app

# Install SQLite runtime dependencies
RUN apk add --no-cache sqlite-libs

COPY --from=build /app/students-api /app/students-api
COPY config/local.yaml /app/config/local.yaml
COPY storage/storage.db /app/storage/storage.db
EXPOSE 8082
CMD ["./students-api", "-config", "config/local.yaml"]
