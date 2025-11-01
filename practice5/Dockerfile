FROM ubuntu:latest
LABEL authors="anaraarman"

ENTRYPOINT ["top", "-b"]
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o products-api .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/products-api .

EXPOSE 8080

CMD ["./products-api"]
