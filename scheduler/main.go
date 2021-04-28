package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"os"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET_3"),
		os.Getenv("LINE_ACCESS_TOKEN_3"),
	)

	if err != nil {
		log.Fatalf("Unable to connect to line bot: %v\n", err)
	}
	user := os.Getenv("LINE_USER_ID")

	unreadMessages, err := process()

	if _, err := bot.PushMessage(user, linebot.NewTextMessage(unreadMessages)).Do(); err != nil {
		log.Println(err)
	}
}
