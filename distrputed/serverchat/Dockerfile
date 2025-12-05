# --- Build stage ---
FROM golang:alpine AS builder
WORKDIR /app

# Enable reproducible builds; speed up module downloads
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Cache modules separately
COPY go.mod ./
RUN go mod download

# Copy source and build
COPY . .
RUN go build -o server .


# --- Final stage ---
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/server /app/server

# If using alpine, add a non-root user (good practice)
RUN adduser -D appuser
USER appuser

EXPOSE 8080
ENTRYPOINT ["/app/server"]

ENV GOTOOLCHAIN=auto