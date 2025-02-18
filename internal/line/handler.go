package line

import (
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	
	"disney-wait-bot/internal/disney"
)

// Webhook のハンドラー（メッセージを受け取る）
func WebhookHandler(req *webhook.CallbackRequest, r *http.Request) {
	// 受信したイベントを処理
	for _, event := range req.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			// テキストメッセージを受け取った場合
			switch e.Message.(type) {
			case webhook.TextMessageContent:
				_, err := LineBot.ReplyMessage(
					&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							&messaging_api.TextMessage{
								Text: "個別返信には対応していません",
							},
						},
					},
				)
				if err != nil {
					log.Println("メッセージの返信に失敗:", err)
				}
			}
		case webhook.PostbackEvent:
			park := e.Postback.Data

			attractionInfoList := disney.GetAttractionInfoList(park)
			var message string

			for _, attractionInfo := range attractionInfoList {
				message += attractionInfo.Name + "\n:      " + attractionInfo.WaitTime + "\n\n"
			}

			_, err := LineBot.ReplyMessage(
				&messaging_api.ReplyMessageRequest{
					ReplyToken: e.ReplyToken,
					Messages: []messaging_api.MessageInterface{
						&messaging_api.TextMessage{
							Text: message,
						},
					},
				},
			)
			if err != nil {
				log.Println("メッセージの返信に失敗:", err)
			}
		default:
			log.Printf("イベントが処理されていません: %T", e)
		}
	}
}
