package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	http.HandleFunc("/", helloHundler)
	http.HandleFunc("/callback", lineHundler)

	fmt.Println("http://localhost:8080で起動中...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHundler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello World!!"
	fmt.Fprintf(w, msg)
}

func lineHundler(w http.ResponseWriter, r *http.Request) {

	bot, err := linebot.New(
        os.Getenv("LINE_BOT_CHANNEL_SECRET"),
        os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
    )
    // エラーに値があればログに出力し終了する
    if err != nil {
        log.Fatal(err)
    }
    // テキストメッセージを生成する
    message := linebot.NewTextMessage("hello, world")
    // テキストメッセージを友達登録しているユーザー全員に配信する
    if _, err := bot.BroadcastMessage(message).Do(); err != nil {
        log.Fatal(err)
    }
	
}