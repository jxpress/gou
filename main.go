package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strings"

	"github.com/nlopes/slack"
	_ "github.com/mattn/go-sqlite3"
)
type Karma struct {
	Giver string
	Receiver string
	Count float64
	Channel string
}

var IgnoreWords = []string {
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

// カルマデータ生成処理
func parseKarma(text string, giver string, channel string) (karmaList []Karma, err error){
	count := float64(strings.Count(text, "+") - 1)  // 雑
	// ++ → 1、 +++ → 1.1、 ++++ → 1.2 ... のように変換する
	if count > 1 {
		count = 1.0 + (count - 1)/ 10.0
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
func giveKarmaEvent(event *slack.MessageEvent) error {
	karmaList, err := parseKarma(event.Text, event.Username, event.Channel)
	if err != nil {
		return err
	}
	fmt.Println("Get karma", karmaList)
	return nil
}

// Botが Join しているチャンネルに投稿されたもの処理する
func handleMessageEvent(event *slack.MessageEvent) error {
	text := event.Text
	if strings.Contains(text, "++") {
		return giveKarmaEvent(event)
	}
	return nil
}

func createKarmaTable(dataDir string) {
	if dataDir == "" {
		dataDir = "."
	}
	db, err := sql.Open("sqlite3", path.Join(dataDir, "karma.db"))
	if err != nil {
		panic(err)
	}
	_, filename, _, _ := runtime.Caller(1)
	raw, err := ioutil.ReadFile(path.Join(path.Dir(filename), "create_table.sql"))
	if err != nil {
		panic(err)
	}
	query := string(raw)
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func main() {

	env := NewEnv()

	createKarmaTable(env.DataDir)

	api := slack.New(
		env.SlackApiKey,
		slack.OptionDebug(env.Debug),
	)

	rtm := api.NewRTM()
	 go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			err := handleMessageEvent(ev)
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
