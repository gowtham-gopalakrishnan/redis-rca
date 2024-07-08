FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o redis-reader

FROM alpine:3.14

RUN apk add --no-cache iproute2

WORKDIR /app

COPY --from=builder /app/redis-reader .

CMD ["./redis-reader"]