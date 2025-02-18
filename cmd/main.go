package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"disney-wait-bot/internal/line"
)

func main() {
	line.InitBot()
	line.CreateRichMenu()
	// Webhook のエンドポイントを設定
	line.LineHandler.HandleEvents(line.WebhookHandler)
	http.Handle("/webhook", line.LineHandler)
	// サーバーを起動
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	fmt.Println("サーバーを起動中: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
