# pakai golang buat build
FROM golang:1.24 AS builder

WORKDIR /app

# copy go mod
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build binary microA
RUN go build -o microa ./services/microA/cmd/microa

# runtime image
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/microa .

CMD ["./microa"]
