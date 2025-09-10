# syntax=docker/dockerfile:1.7
FROM golang:1.24 AS builder
WORKDIR /app

# Pastikan git & CA ada (untuk fetch modul via VCS)
RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && \
    apt-get install -y --no-install-recommends git ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Pakai proxy resmi + fallback ke direct (hindari goproxy.io)
ENV GOPROXY=https://proxy.golang.org,direct \
    GOSUMDB=sum.golang.org

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o microb ./services/microB/cmd/microb

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/microb .
CMD ["./microb"]
