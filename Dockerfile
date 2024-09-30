# Building stage
FROM golang:1.21.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build


# Running stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/pyre-promotion .
COPY --from=builder /app/.env .

EXPOSE 8000

CMD ["./pyre-promotion"]