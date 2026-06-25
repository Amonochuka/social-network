# --------------------------------------------------
# Build stage
# --------------------------------------------------
FROM golang:1.25-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build \
    -a \
    -ldflags="-s -w" \
    -o social-api \
    ./backend/cmd/api

# --------------------------------------------------
# Runtime stage
# --------------------------------------------------
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /src/social-api .

# SQLite database location
COPY app.db ./app.db

# Upload directories required by the application
RUN mkdir -p \
    uploads/avatars \
    uploads/comments \
    uploads/posts

EXPOSE 8080

CMD ["./social-api"]