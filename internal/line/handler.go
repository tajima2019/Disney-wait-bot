package line

import (
	"encoding/json"
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
		case webhook.MessageEvent: // メッセージイベントの場合
			switch e.Message.(type) {
			case webhook.TextMessageContent: // テキストメッセージの場合
				text := e.Message.(webhook.TextMessageContent).Text
				attractionInfoListByInitial := disney.GetAttractionInfoListByInitial(text)
				contents := createContents(attractionInfoListByInitial)
				flexMessage := CreateFlexMessage(contents)
				_, err := LineBot.ReplyMessage(
					&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							&messaging_api.FlexMessage{
								AltText: text + "から始まるアトラクションの待ち時間",
								Contents: flexMessage,
							},
						},
					},
				)
				if err != nil {
					log.Println("メッセージの返信に失敗:", err)
				}
			}
		case webhook.PostbackEvent: // ポストバックイベントの場合
			park := e.Postback.Data

			attractionInfoList := disney.GetAttractionInfoList(disney.Park(park))

			contents := createContents(attractionInfoList)

			flexMessage := CreateFlexMessage(contents)

			// `Flex Message` を送信
			_, err := LineBot.ReplyMessage(
				&messaging_api.ReplyMessageRequest{
					ReplyToken: e.ReplyToken,
					Messages: []messaging_api.MessageInterface{
						&messaging_api.FlexMessage{
							AltText:  "ディズニーの待ち時間",
							Contents: flexMessage,
						},
					},
				},
			)
			if err != nil {
				log.Println("メッセージの返信に失敗:", err)
			}
		default: // その他のイベントの場合
			log.Printf("イベントが処理されていません: %T", e)
		}
	}
}

func createContents(attractionInfoList []disney.AttractionInfo) (contents []interface{}) {
	contents = []interface{}{
		map[string]interface{}{
			"type":   "text",
			"text":   "🎢 ディズニーの待ち時間",
			"weight": "bold",
			"size":   "xl",
			"color":  "#1DB446", // 緑色でタイトルを強調
		},
		map[string]interface{}{
			"type":   "separator",
			"margin": "md",
		},
	}
	// ✅ Attraction データを追加（無制限に増やせる）
	for _, attractionInfo := range attractionInfoList {
		contents = append(contents, map[string]interface{}{
			"type": "box",
			"layout": "vertical",
			"contents": []interface{}{
				// アトラクションの名前
				map[string]interface{}{
					"type":   "text",
					"text":   attractionInfo.Name,
					"weight": "bold",
					"size":   "md",
					"color":  "#FF0000",
				},
				// 待ち時間
				map[string]interface{}{
					"type":   "text",
					"text":   "⏳ " + attractionInfo.WaitTime,
					"size":   "md",
					"color":  "#333333",
					"weight": "bold",
					"margin": "sm",
				},
			},
			"paddingAll": "10px",
			"backgroundColor": "#F0F0F0",
			"cornerRadius": "8px",
		})
		// ✅ 各アトラクションの間に区切り線を入れる（ただし、最後には入れない）
		if attractionInfo != attractionInfoList[len(attractionInfoList)-1] {
			contents = append(contents, map[string]interface{}{
				"type":   "separator",
				"margin": "md",
			})
		}
	}
	return contents
} 

func CreateFlexMessage(contents []interface{}) (flexMessage messaging_api.FlexContainerInterface){
	// `Flex Message` の JSON を組み立て
	flexMessageJSON := map[string]interface{}{
		"type": "bubble",
		"body": map[string]interface{}{
			"type":     "box",
			"layout":   "vertical",
			"contents": contents,
		},
	}

	// JSON を変換
	flexMessageBytes, _ := json.Marshal(flexMessageJSON)
	flexMessage, err := messaging_api.UnmarshalFlexContainer(flexMessageBytes)
	if err != nil {
		log.Println("Flex Message の生成に失敗:", err)
		return
	}
	return flexMessage
} 