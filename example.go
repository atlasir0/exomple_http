package main

import (
	"fmt"
	"net/http"
	"strings"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const (
	BotToken   = "7145361114:AAGcDmLWHv9eyeQTjcj1djRA1oDcCJBmuKg"
	WebhookURL = "https://9fb9-178-207-154-253.ngrok-free.app"
)

var (
	messageMap = make(map[string]string)
	imageMap   = make(map[string]string)
)

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
	command := strings.ToLower(update.Message.Text)
	if response, ok := messageMap[command]; ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		bot.Send(msg)
	} else if imageURL, ok := imageMap[command]; ok {
		msg := tgbotapi.NewPhotoShare(update.Message.Chat.ID, imageURL)
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Command not found")
		bot.Send(msg)
	}
}

func start() {
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
	messageMap["/help"] = "Available commands:\n/help - Show available commands\n/add - Add a new command\n"
	messageMap["/add"] = " ins  '/add hello:Hi there!'"
	messageMap["/es"] = "Wtf"
	imageMap["/cat"] = "https://s3.amazonaws.com/freecodecamp/running-cats.jpg"

	start()
}
