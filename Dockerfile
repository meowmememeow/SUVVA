FROM golang:1.22.3 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /suvva-geo-ride-service

FROM alpine:latest
COPY --from=builder /suvva-geo-ride-service /suvva-geo-ride-service
RUN chmod +x /suvva-geo-ride-service
CMD ["/suvva-geo-ride-service"]
