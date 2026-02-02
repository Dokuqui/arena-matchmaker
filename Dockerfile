FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o frontend cmd/frontend/*.go
RUN go build -o matchmaker cmd/matchmaker/*.go

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/frontend .
COPY --from=builder /app/matchmaker .
COPY --from=builder /app/web ./web

EXPOSE 50051 2112 8080

CMD ["./frontend"]
