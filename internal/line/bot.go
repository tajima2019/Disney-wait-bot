package line

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

var (
	LineBot       *messaging_api.MessagingApiAPI
	LineBlobClient *messaging_api.MessagingApiBlobAPI
	LineHandler *webhook.WebhookHandler
)

func InitBot() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("環境変数をロードできませんでした")
	}

	LineBot, err = messaging_api.NewMessagingApiAPI(os.Getenv("LINE_CHANNEL_TOKEN"))
	if err != nil {
		log.Fatal("LINE Bot の Bot 初期化に失敗:", err)
	}

	LineBlobClient, err = messaging_api.NewMessagingApiBlobAPI(os.Getenv("LINE_CHANNEL_TOKEN"))
	if err != nil {
		log.Fatal("LINE Bot の Blob_client 初期化に失敗:", err)
	}

	LineHandler, err = webhook.NewWebhookHandler(os.Getenv("LINE_CHANNEL_SECRET"))
	if err != nil {
		log.Fatal("LINE Bot の Handler 初期化に失敗:", err)
	}
}
