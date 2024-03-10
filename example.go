package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const (
	BotToken   = "7145361114:AAGcDmLWHv9eyeQTjcj1djRA1oDcCJBmuKg"
	WebhookURL = "https://87af-178-207-154-253.ngrok-free.app"
)

var rss = map[string]string{
	"Habr": "https://habrahabr.ru/rss/best/",
}

type RSS struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	URL   string `xml:"guid"`
	Title string `xml:"title"`
}

func initializeBot() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	return bot, nil
}

func setWebhook(bot *tgbotapi.BotAPI) error {
	_, err := bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	return err
}

func startListening(bot *tgbotapi.BotAPI) {
	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(":8080", nil)
	fmt.Println("start listen :8080")

	for update := range updates {
		handleMessage(bot, update)
	}
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Text {
	case "/help":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Available commands:\n"+
			"/help - Show available commands\n"+
			"Enter 'Habr' to get news from Habr")
		bot.Send(msg)
	case "/add":
		add := tgbotapi.NewMessage(update.Message.Chat.ID, "if you want add new file enter:\n"+"/add")
		bot.Send(add)
	default:
		if url, ok := rss[update.Message.Text]; ok {
			rss, err := getNews(url)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"sorry, error happened",
				))
			}
			for _, item := range rss.Items {
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					item.URL+"\n"+item.Title,
				))
			}
		} else {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"there is only Habr feed available",
			))
		}
	}
}

func getNews(url string) (*RSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rss := new(RSS)
	err = xml.Unmarshal(body, rss)
	if err != nil {
		return nil, err
	}

	return rss, nil
}

func dop() {
	bot, err := initializeBot()
	if err != nil {
		panic(err)
	}

	err = setWebhook(bot)
	if err != nil {
		panic(err)
	}

	startListening(bot)
}

func main() {
	dop()
}
