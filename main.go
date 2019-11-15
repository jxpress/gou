package main

import (
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nlopes/slack"
)

var IgnoreWords = []string{
	"\n",
	"`",
	"　",
}

func containIgnoreWords(target string) bool {
	for _, word := range IgnoreWords {
		if strings.Contains(target, word) {
			return true
		}
	}
	return false
}

type KarmaBot struct {
	repo KarmaRepo
}


// カルマデータ生成処理
func parseKarma(text string, giver string, channel string) (karmaList []Karma, err error) {
	count := float64(strings.Count(text, "+") - 1) // 雑
	// ++ → 1、 +++ → 1.1、 ++++ → 1.2 ... のように変換する
	if count > 1 {
		count = 1.0 + (count-1)/10.0
	}

	// + の前のテキストを抽出する
	sp := strings.Split(text, "+")
	receivers := sp[0]
	for _, name := range strings.Split(receivers, " ") {
		if containIgnoreWords(name) {
			continue
		}
		name = strings.ReplaceAll(name, "@", "")
		name = strings.TrimSpace(name)
		karmaList = append(karmaList, Karma{
			Giver:    giver,
			Receiver: name,
			Count:    count,
			Channel:  channel,
		})
	}
	return
}

// カルマ付与イベント
func (k *KarmaBot) giveKarmaEvent(slack *slack.Client, event *slack.MessageEvent) error {
	karmaList, err := parseKarma(event.Text, event.Username, event.Channel)
	if err != nil {
		return err
	}
	err = k.repo.Save(karmaList)
	fmt.Println("Get karma", karmaList)

	return err
}

// Botが Join しているチャンネルに投稿されたもの処理する
func handleMessageEvent(bot *KarmaBot, slack *slack.Client, event *slack.MessageEvent) error {
	text := event.Text
	if strings.Contains(text, "++") {
		return bot.giveKarmaEvent(slack, event)
	}
	return nil
}

func main() {

	env := NewEnv()

	repo := NewSQLiteKarmaRepo(env.DataDir)
	bot := &KarmaBot{repo:repo}
	api := slack.New(
		env.SlackApiKey,
		slack.OptionDebug(env.Debug),
	)

	fmt.Println("Start receiving...")
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			err := handleMessageEvent(bot, api, ev)
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
