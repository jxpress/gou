package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nlopes/slack"
)

func initUser(repo UserRepo, client *slack.Client) error {
	us, err := client.GetUsers()
	if err != nil {
		return err
	}
	for _, u := range us {
		if u.IsBot {
			continue
		}
		if u.Deleted {
			continue
		}

		err := repo.Put(User{
			Id:          u.ID,
			Name:        u.Name,
			DisplayName: u.RealName,
			TeamId:      u.TeamID,
			ImageURL:    u.Profile.Image192,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {

	env := NewEnv()

	repo := NewSQLiteKarmaRepo(env.DataDir)

	client := slack.New(
		env.SlackApiKey,
		slack.OptionDebug(env.Debug),
	)

	bot, err := NewKarmaBot(repo, client)
	if err != nil {
		panic(err)
	}

	fmt.Println("Start receiving...")
	rtm := client.NewRTM()

	go initUser(repo, client)

	go rtm.ManageConnection()
	for msg := range rtm.IncomingEvents {
		err := eventReceiver(bot, msg)
		if err != nil {
			fmt.Printf("Error %v", err.Error())
		}
	}
}
