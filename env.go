package main

import "os"

type Env struct {
	SlackApiKey string  // Slack のボットのトークン
	Debug bool  // デバッグログを表示するかどうか
	DataDir string  // SQLite のデータを設置するディレクトリ
}


func NewEnv() Env {
	return Env{
		SlackApiKey:os.Getenv("SLACK_API_KEY"),
		Debug: os.Getenv("DEBUG") == "1",
		DataDir: os.Getenv("DATA_DIR"),
	}
}
