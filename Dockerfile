# 1. ビルド用のステージ
FROM golang:1.23 AS builder
WORKDIR /app

# 2. Goのモジュールを設定
COPY go.mod go.sum ./
RUN go mod download

# 3. ソースコードをコピーしてビルド
COPY . .
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# 4. 実行用の軽量コンテナ
FROM alpine:latest
WORKDIR /app

# 5. 必要なランタイムをインストール
RUN apk --no-cache add ca-certificates

# 6. ビルドしたバイナリをコピー
COPY --from=builder /app/main .
COPY --from=builder /app/assets ./assets
COPY .env .

# 7. ポートを指定（必要なら）
EXPOSE 8080

# 8. 実行
ENTRYPOINT ["./main"]
