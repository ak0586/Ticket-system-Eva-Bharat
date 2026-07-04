# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ticket-system .

# Run stage — scratch or distroless keeps the image tiny
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/ticket-system .

EXPOSE 8080

CMD ["./ticket-system"]
