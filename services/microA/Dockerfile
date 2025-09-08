FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /microa ./cmd/microa

FROM alpine:3.18
COPY --from=build /microa /microa
EXPOSE 8090
ENTRYPOINT ["/microa"]
