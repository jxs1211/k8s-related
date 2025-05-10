# Build stage
FROM docker.io/golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./server

# Final stage
FROM docker.io/alpine:latest
WORKDIR /app
COPY --from=builder /server /app/server
EXPOSE 50051
ENTRYPOINT ["/app/server"]
