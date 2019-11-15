package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"strings"
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
	slack *slack.Client
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
		if name == "" {
			continue
		}
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
func (k *KarmaBot) giveKarmaEvent(event *slack.MessageEvent) error {

	karmaList, err := parseKarma(event.Text, event.Username, event.Channel)
	if err != nil {
		return err
	}

	err = k.repo.Save(karmaList)
	if err != nil {
		return err
	}

	fmt.Println("Get karma", karmaList)
	// リアクションをつける
	ref := slack.ItemRef{Timestamp: event.Timestamp, Channel: event.Channel}
	err = k.slack.AddReaction("thumbsup_all", ref)
	return err
}

// Botが Join しているチャンネルに投稿されたもの処理する
func handleMessageEvent(bot *KarmaBot, event *slack.MessageEvent) error {
	text := event.Text
	if strings.Contains(text, "++") {
		return bot.giveKarmaEvent(event)
	}
	return nil
}

func eventReceiver(bot *KarmaBot, msg slack.RTMEvent) error {
	switch ev := msg.Data.(type) {
	case *slack.MessageEvent:
		err := handleMessageEvent(bot, ev)
		if err != nil {
			return err
		}

	case *slack.RTMError:
		return fmt.Errorf("rtm error: %v", ev)

	case *slack.InvalidAuthEvent:
		return fmt.Errorf("invalid credentials: %v", ev)
	default:
	}
	return nil
}
