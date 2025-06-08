# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o backend main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/backend .
COPY --from=builder /app/model ./model
COPY wait-for-postgres.sh /wait-for-postgres.sh
RUN chmod +x /wait-for-postgres.sh
ENV DB_HOST=postgres
EXPOSE 8080
CMD ["/bin/sh", "-c", "/wait-for-postgres.sh $DB_HOST:5432 60 && ./backend"]
