package main

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"net/http"
	"os"
	"time"
)

type info struct {
	Date    string
	From    string
	Subject string
}

func connectToGmail() *http.Client {
	config := oauth2.Config{
		ClientID:     os.Getenv("GMAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("GMAIL_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.readonly"},
	}

	expiry, _ := time.Parse("2006-01-02T15:04:05+09:00 (MST)", "2020-08-17T12:13:44.918393+09:00 (JST)")
	token := oauth2.Token{
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
		TokenType:    "Bearer",
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		Expiry:       expiry,
	}
	client := config.Client(oauth2.NoContext, &token)

	return client
}

func process() (string, error) {
	var unreadMessages []info

	client := connectToGmail() // generate client

	srv, err := gmail.New(client)
	if err != nil {
		return "", err
	}

	user := os.Getenv("MAIL")
	msgs, err := srv.Users.Messages.List(user).Q("is:unread").Do()
	if err != nil {
		return "", err
	}

	for _, msgMsg := range msgs.Messages {
		id := msgMsg.Id                                   // id of a message
		msg, err := srv.Users.Messages.Get(user, id).Do() // get message by its id
		if err != nil {
			return "", err
		}
		var date, from, subject string
		// collect date, from, subject from each message
		for _, header := range msg.Payload.Headers {
			if header.Name == "Date" {
				date = header.Value
			} else if header.Name == "From" {
				from = header.Value
			} else if header.Name == "Subject" {
				subject = header.Value
			}
		}
		// make sure to complete three fields
		if date == "" || from == "" || subject == "" {
			continue
		}
		msgInfo := info{Date: date, From: from, Subject: subject}
		unreadMessages = append(unreadMessages, msgInfo)
	}

	res, err := toString(unreadMessages)
	if err != nil {
		return "", err
	}
	return res, nil
}

func toString(unreadMessages []info) (string, error) {
	res := ""
	for i, msg := range unreadMessages {
		// extract year/month/day from msg.Date
		date, err := time.Parse(time.RFC1123Z, msg.Date)
		formattedDate := date.Format("2006/01/02")
		if err != nil {
			return "", err
		}

		// extract sender's mail address from msg.From
		mail := msg.From

		// extract subject from msg.Subject
		sub := msg.Subject

		if i == len(unreadMessages)-1 {
			res += fmt.Sprintf("%s\n(%s)\n[%s]", formattedDate, mail, sub)
		} else {
			res += fmt.Sprintf("%s\n(%s)\n[%s]\n\n", formattedDate, mail, sub)
		}
	}
	return res, nil
}
