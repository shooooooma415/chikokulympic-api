FROM golang:1.24.3-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache git sqlite-dev gcc musl-dev

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

# Cloud Runは8080ポートを期待する
ENV PORT=8080
EXPOSE 8080

CMD ["./main"]