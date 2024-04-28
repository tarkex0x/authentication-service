FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]