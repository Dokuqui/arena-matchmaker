FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server cmd/server/*.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 50051

CMD ["./server"]
