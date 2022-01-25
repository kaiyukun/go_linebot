package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"unicode/utf8"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	
	
)

func main() {
	// ハンドラの登録
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/callback", lineHandler)
	
	fmt.Println("http://localhost:8080 で起動中...")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
        log.Fatal(err)
    }
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	msg := "Hello World!!!!"
	fmt.Fprintf(w, msg)
}

func lineHandler(w http.ResponseWriter, r *http.Request) {
	// BOTを初期化
	bot, err := linebot.New(
		// os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		// os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
		"63807e1aa6417dd7b78c01347f50b473",
		"7iU4LtyQA8sruQboSNuFIgHrOo0CrgpR20TKDH0nq6ibOo0JBUUMu7SZ3mZ5l01oPlJ+U8BY2vPtOyZiO4zAQeEw/FnQO6vsDCkFq7/zUuUf5sIh+fRUI+VAZPPJnbNktykFPajdPzCJXihmXPOofwdB04t89/1O/w1cDnyilFU=",
	)
	if err != nil {
		log.Fatal(err)
	}

	// リクエストからBOTのイベントを取得
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		// イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// メッセージがテキスト形式の場合
			case *linebot.TextMessage:
				replyMessage := message.Text
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
				if err != nil {
					log.Print(err)
				}
			case *linebot.LocationMessage:
				sendRestoInfo(bot, event)
			}

		}
	}
}

func sendRestoInfo(bot *linebot.Client, e *linebot.Event) {
	msg := e.Message.(*linebot.LocationMessage)

	lat := strconv.FormatFloat(msg.Latitude, 'f', 2, 64)
	lng := strconv.FormatFloat(msg.Longitude, 'f', 2, 64)

	replyMsg := getRestoInfo(lat, lng)

	res := linebot.NewTemplateMessage(
		"レストラン一覧",
		linebot.NewCarouselTemplate(replyMsg...).WithImageOptions("rectangle", "cover"),
	)
	if _, err := bot.ReplyMessage(e.ReplyToken, res).Do(); err != nil {
		log.Print(err)
	}
}

// response APIレスポンス
type response struct {
	Results results `json:"results"`
}

// results APIレスポンスの内容
type results struct {
	Shop []shop `json:"shop"`
}

// shop レストラン一覧
type shop struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Open    string `json:"open"`
	Photo photo `json:"photo"`
	URLS urls `json:"urls"`
}

type photo struct {
	Mobile mobile `json:"mobile"`
}

type mobile struct {
	L string `json:"l"`	
}

type urls struct {
	PC string `json:"pc"`	
}

func getRestoInfo(lat string, lng string) []*linebot.CarouselColumn {
	apikey := "4869eef30bfe5c4a"
	url := fmt.Sprintf(
		"https://webservice.recruit.co.jp/hotpepper/gourmet/v1/?format=json&key=%s&lat=%s&lng=%s&range=5",
		apikey, lat, lng)

	// リクエストしてボディを取得
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data response
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	var ccs []*linebot.CarouselColumn
	for _, shop := range data.Results.Shop {
		addr := shop.Address
		if 60 < utf8.RuneCountInString(addr) {
			addr = string([]rune(addr)[:60])
		}

		cc := linebot.NewCarouselColumn(
			shop.Photo.Mobile.L,
			shop.Name + shop.Open,
			addr,
			linebot.NewURIAction("ホットペッパーで開く", shop.URLS.PC),
		).WithImageOptions("#FFFFFF")
		ccs = append(ccs, cc)
	}
	return ccs
}