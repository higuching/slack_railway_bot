package app

import (
	"bytes"
	"log"
	"net/http"

	db "github.com/higuching/slack_bot/db"
)

type SlackJson struct {
	Token string
}

var botId string

func Run() (string, error) {
	err := db.NewRailways().Create()
	if err != nil {
		return "", err
	}
	log.Print("Database created.")
	name := "RailwaysBot"
	text := NewRailWays().GetMessage()
	channel := "#bot用"
	if text == "" {
		return "", nil
	}

	//jsonStr := `{"username":"` + name + `","text":"` + text + `"}`
	jsonStr := `{"channel":"` + channel + `","username":"` + name + `","text":"` + text + `"}`

	req, err := http.NewRequest(
		"POST",
		"https://hooks.slack.com/services/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",
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		return text, err
	}

	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return text, err
	}
	defer resp.Body.Close()

	return text, nil

}
