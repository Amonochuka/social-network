# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

# CGO is required for go-sqlite3
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Cache module downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -o /server ./backend/cmd/api

# ── Stage 2: Runtime ──────────────────────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /server .
COPY --from=builder /app/backend/pkg/db/migrations ./backend/pkg/db/migrations

# Create upload directory
RUN mkdir -p uploads/avatars

EXPOSE 8080

CMD ["./server"]
