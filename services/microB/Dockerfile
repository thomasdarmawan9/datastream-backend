FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /microb ./cmd/microb

FROM alpine:3.18
COPY --from=build /microb /microb
EXPOSE 8080
ENTRYPOINT ["/microb"]
