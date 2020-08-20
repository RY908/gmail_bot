package main

import (
	"fmt"
	"log"
	"os"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)


func main() {
	// connect to Line bot
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET_3"), 
		os.Getenv("LINE_ACCESS_TOKEN_3"),
	)

	if err != nil {
		fmt.Println(err)
	}
	
	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
				if err == linebot.ErrInvalidSignature {
						w.WriteHeader(400)
				} else {
						w.WriteHeader(500)
				}
				return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
					switch message := event.Message.(type) {
					case *linebot.TextMessage:
							if event.ReplyToken == "00000000000000000000000000000000" {
									return
							}
							// if the message user send is "メールを確認", then bot send unread messages.
							if message.Text == "メールを確認" {
								// search unread messages
								unreadMessages, err := process()
								if err != nil {
									log.Fatalf("Unable to retrieve unread messages: %v", err)
								}
								if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(unreadMessages)).Do(); err != nil {
										log.Print(err)
								}
							} else {
								replayMessage := "メールを確認するには「メールを確認」と送信してください。"
								if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replayMessage)).Do(); err != nil {
									log.Print(err)
								}
							}
						
					}
			}
		}
	})
	// This is just sample code.
  // For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
  if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
	

}

//https://github.com/heroku/heroku-buildpack-go.git