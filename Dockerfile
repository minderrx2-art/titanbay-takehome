FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Use CGO_ENABLED=0 to create a statically linked binary (portable)
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Alpine runner makes the image tiny (~20MB vs ~200MB for Ubuntu)
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]