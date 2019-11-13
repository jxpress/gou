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
