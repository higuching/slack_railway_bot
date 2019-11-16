package main

import (
	"fmt"
	"github.com/higuching/slack_railway_bot/app"
	"os"
)

func main() {
	_, err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(0)
}
