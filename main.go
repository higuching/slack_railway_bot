package main

import (
	"github.com/higuching/slack_railway_bot/app"
)

func main() {
	_, err := app.Run()
	if err != nil {
		panic(err)
	}
}
