# STAGE 1
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# STAGE 2
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

ENTRYPOINT ["./main"]
