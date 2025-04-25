# ビルドステージ
FROM golang:1.23-bullseye AS builder
WORKDIR /app

# 依存解決
COPY go.mod go.sum ./
RUN go mod download

# ソースをコピーしてビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd


# 実行ステージ
FROM gcr.io/distroless/base
WORKDIR /cmd
COPY --from=builder /app/server .

# Echo はデフォルトで 8080 ポートを使う想定
EXPOSE 8080
CMD ["./server"]