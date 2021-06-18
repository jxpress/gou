package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"strings"
	"time"
)

var ValidEmojis = []string{
	"+1",
	"god",
	"sbrr",
}
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
	userKarmaRepo UserKarmaRepo
	slack         *slack.Client
	userId        string  // Bot 自身のユーザーID
}

func NewKarmaBot (repo UserKarmaRepo, client *slack.Client) (*KarmaBot, error) {
	bot := &KarmaBot{userKarmaRepo: repo, slack: client}
	// API を叩いてBot自身のIDを取得
	res, err := client.AuthTest()
	if err != nil {
		return nil, err
	}
	bot.userId = res.UserID
	return bot, nil
}

func (k *KarmaBot) parseUser(nameOrId string) User {
	user, err := k.userKarmaRepo.GetById(nameOrId)
	if err == nil {
		return user
	}
	user, _ = k.userKarmaRepo.GetByName(nameOrId)
	return user
}

// カルマデータ生成処理
func (k *KarmaBot) parseKarma(text string, giver string, channel string) (karmaList []Karma, err error) {
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
		name = strings.ReplaceAll(name, "<", "")
		name = strings.ReplaceAll(name, ">", "")
		name = strings.TrimSpace(name)
		user := k.parseUser(name)
		karmaList = append(karmaList, Karma{
			Giver:    giver,
			Receiver: user.Id,
			Count:    count,
			Channel:  channel,
		})
	}
	return
}

// カルマ付与イベント
func (k *KarmaBot) giveKarmaEvent(event *slack.MessageEvent) error {

	karmaList, err := k.parseKarma(event.Text, event.User, event.Channel)
	if err != nil {
		return err
	}

	err = k.userKarmaRepo.Save(karmaList)
	if err != nil {
		return err
	}

	fmt.Println("Get karma", karmaList)
	// リアクションをつける
	ref := slack.ItemRef{Timestamp: event.Timestamp, Channel: event.Channel}
	err = k.slack.AddReaction("thumbsup_all", ref)
	return err
}

// ランキング表示イベント
func (k *KarmaBot) getKarmaRankingEvent(event *slack.MessageEvent) error {
	from := time.Now().AddDate(0, 0, -7)
	to := time.Now()
	ranking, err := k.userKarmaRepo.Ranking(KindReceiver, from, to)
	if err != nil {
		return err
	}
	msg := ""
	for i, r := range ranking.Ranks {
		msg += fmt.Sprintf("%d. <@%s> (%.2fpt)\n", i + 1, r.User, r.Count)
	}
	text := slack.MsgOptionText(msg, false)
	ts := slack.MsgOptionTS(event.Timestamp)
	emoji := slack.MsgOptionIconEmoji(":karma:")
	_, _, err = k.slack.PostMessage(event.Channel, text, ts, emoji)
	return err
}

func (k *KarmaBot) showHelpEvent(event *slack.MessageEvent) error {
	text := slack.MsgOptionText(`ランキングを見たい場合はメンションのあとに「ランキング」とつけてください。1週間分のランキングがみれます`, false)
	ts := slack.MsgOptionTS(event.Timestamp)
	emoji := slack.MsgOptionIconEmoji(":karma:")
	_, _, err := k.slack.PostMessage(event.Channel, text, ts, emoji)
	return err
}

// Botが Join しているチャンネルに投稿されたもの処理する
func handleMessageEvent(bot *KarmaBot, event *slack.MessageEvent) error {
	text := event.Text
	if strings.Contains(text, "++") {
		return bot.giveKarmaEvent(event)
	}

	// Bot に対するリプライ
	if strings.Contains(text, bot.userId) {
		if strings.Contains(text, "ランキング") {
			return bot.getKarmaRankingEvent(event)
		}
		return bot.showHelpEvent(event)
	}
	return nil
}

// スタンプに対する処理をする
func handleReactionEvent(bot *KarmaBot, event *slack.ReactionAddedEvent) error {
	var ok bool
	for _, w := range ValidEmojis {
		if strings.HasPrefix(event.Reaction, w) {
			ok = true
		}
	}
	if !ok {
		return nil
	}
	karma := Karma{
		Giver:    event.User,
		Receiver: event.ItemUser,
		Count:    0.1,
		Channel:  event.Item.Channel,
	}
	if err := bot.userKarmaRepo.Save([]Karma{karma}); err != nil {
		return err
	}
	// リアクションをつける
	ref := slack.ItemRef{Timestamp: event.Item.Timestamp, Channel: event.Item.Channel}
	err := bot.slack.AddReaction("eyes", ref)
	return err
}

func eventReceiver(bot *KarmaBot, msg slack.RTMEvent) error {
	switch ev := msg.Data.(type) {
	case *slack.ReactionAddedEvent:
		err := handleReactionEvent(bot, ev)
		if err != nil {
			return err
		}
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
