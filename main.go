package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nlopes/slack"
)

func main() {

	env := NewEnv()

	repo := NewSQLiteKarmaRepo(env.DataDir)

	client := slack.New(
		env.SlackApiKey,
		slack.OptionDebug(env.Debug),
	)
	bot,err := NewKarmaBot(repo, client)

	if err != nil {
		panic(err)
	}

	fmt.Println("Start receiving...")
	rtm := client.NewRTM()

	go rtm.ManageConnection()
	for msg := range rtm.IncomingEvents {
		err := eventReceiver(bot, msg)
		if err != nil {
			fmt.Printf("Error %v", err.Error())
		}
	}
}
