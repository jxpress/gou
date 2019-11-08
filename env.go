package main

import "os"

type Env struct {
	SlackApiKey string
	Debug bool
}


func NewEnv() Env {
	return Env{
		SlackApiKey:os.Getenv("SLACK_API_KEY"),
		Debug: os.Getenv("DEBUG") == "1",
	}
}
