package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jxpress/gou/infrastructure"
	"github.com/jxpress/gou/karma"
	"strings"

	"github.com/nlopes/slack"
)

func parseKarma(text string) (users []string, count int, err error){
	sp := strings.Split(text, "+")
	names := sp[0]
	for _, name := range strings.Split(names, " ") {
		name = strings.ReplaceAll(name, "@", "")
		name = strings.TrimSpace(name)
		users = append(users, name)
	}
	count = strings.Count(text, "+")  // 雑
	return

}
func handleEvent(event *slack.MessageEvent) error {
	text := event.Text
	if ! strings.Contains(text, "++") {
		return nil
	}

	users, count, err := parseKarma(text)
	if err != nil {
		return err
	}
	fmt.Println("Get karma", users, count)
	var putErr error
	for _, user := range users {
		putErr = karma.Put(&aws.Config{Region: aws.String("ap-north-east-1")}, &infrastructure.Karma{
			Identifier: user,
			Score:      count,
		})
	}
	return putErr
}

func main() {
	env := NewEnv()
	api := slack.New(
		env.SlackApiKey,
		slack.OptionDebug(env.Debug),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			err := handleEvent(ev)
			if err != nil {
				fmt.Printf("Error %v", err.Error())
			}

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return
		default:
		}
	}
}
