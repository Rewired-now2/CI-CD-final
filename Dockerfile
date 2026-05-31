FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o tracker .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/tracker .

CMD ["./tracker"]