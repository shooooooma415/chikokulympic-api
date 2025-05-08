FROM golang:1.20.2-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache git sqlite-dev gcc musl-dev

COPY . .
RUN go build -o main ./cmd/main.go


FROM alpine
RUN apk update && apk upgrade
RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/main ./main
COPY --from=builder /app/config ./config
COPY --from=builder /app/domain ./domain
COPY --from=builder /app/infrastructure ./infrastructure
COPY --from=builder /app/middleware ./middleware
COPY --from=builder /app/presentation ./presentation
COPY --from=builder /app/server ./server
COPY --from=builder /app/usecase ./usecase
COPY --from=builder /app/.env* ./ 

CMD ["./main"]