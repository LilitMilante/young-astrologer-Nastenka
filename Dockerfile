FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/astrologer

FROM alpine
WORKDIR /app
COPY --from=builder /app/astrologer /app/astrologer
EXPOSE 80
CMD ["/app/astrologer"]

