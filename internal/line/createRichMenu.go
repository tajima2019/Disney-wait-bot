package line

import (
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

// リッチメニューを作成する関数
func CreateRichMenu() {
	// リッチメニューの定義
	richMenu := &messaging_api.RichMenuRequest{
		Size:        &messaging_api.RichMenuSize{Width: 2500, Height: 843},
		Selected:    true,
		Name:        "ディズニー待ち時間メニュー",
		ChatBarText: "メニューを開く",
		Areas: []messaging_api.RichMenuArea{
			{
				Bounds: &messaging_api.RichMenuBounds{X: 0, Y: 0, Width: 1250, Height: 421},
				Action: &messaging_api.PostbackAction{
					Data: "land",
				},
			},
			{
				Bounds: &messaging_api.RichMenuBounds{X: 1250, Y: 0, Width: 1250, Height: 421},
				Action: &messaging_api.PostbackAction{
					Data: "sea",
				},
			},
		},
	}

	res, err := LineBot.CreateRichMenu(richMenu)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open("assets/rich_menu.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	contentType := "image/png"
	richMenuId := res.RichMenuId
	if _, err := LineBlobClient.SetRichMenuImage(richMenuId, contentType, file); err != nil {
		log.Println("リッチメニューの画像の設定に失敗:", err)
		log.Fatal(err)
	}
	
	if _, err := LineBot.SetDefaultRichMenu(richMenuId); err != nil {
		log.Println("リッチメニューの設定に失敗:", err)
		log.Fatal(err)
	}
}
