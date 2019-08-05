package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	db "github.com/higuching/slack_bot/db"
	"github.com/nlopes/slack"
)

type SlackJson struct {
	Token string
}

var botId string

func Run() {
	err := db.NewRailways().Create()
	if err != nil {
		panic(err)
	}
	log.Print("Database created.")
	client := &http.Client{}
	for {
		name := "RailwaysBot"
		text := NewRailWays().GetMessage()
		channel := "#random"
		if text != "" {

			//jsonStr := `{"username":"` + name + `","text":"` + text + `"}`
			jsonStr := `{"channel":"` + channel + `","username":"` + name + `","text":"` + text + `"}`

			req, err := http.NewRequest(
				"POST",
				"https://hooks.slack.com/services/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",
				bytes.NewBuffer([]byte(jsonStr)),
			)
			if err != nil {
				panic(err)
			}

			req.Header.Set("content-type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

		}
		time.Sleep(60 * time.Second)
	}
}

func Run2() {
	token := getToken()
	if token == "" {
		log.Fatal("void token")
	}
	api := slack.New(token)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				err := db.NewRailways().Create()
				if err != nil {
					panic(err)
				}
				log.Print("Database created.")

			case *slack.ConnectedEvent:
				botId = ev.Info.User.ID

			case *slack.MessageEvent:
				if strings.HasPrefix(ev.Text, "<@"+botId+">") {
					text := getText(ev.Text)
					switch text {
					case "遅延":
						message := NewRailWays().GetMessage()
						rtm.SendMessage(rtm.NewOutgoingMessage(message, ev.Channel))
					case "help":
						rtm.SendMessage(rtm.NewOutgoingMessage("Usage: 運行状況　関東圏の電車の運行状況を表示します。\n他の機能？ないよ。", ev.Channel))
					default:
						// rtm.SendMessage(rtm.NewOutgoingMessage("登録されていないパターンです", ev.Channel))
					}
				}

			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				os.Exit(1)

			}
		}
	}
}

// Jsonに記述したbotのTokenを取得する
func getToken() string {
	// JSONファイル読み込み
	bytes, err := ioutil.ReadFile("configs/slack.json")
	if err != nil {
		log.Fatal(err)
	}
	// JSONデコード
	var slackJson SlackJson
	if err := json.Unmarshal(bytes, &slackJson); err != nil {
		log.Fatal(err)
	}
	return slackJson.Token
}

// 送られてきたメッセージのIDを除いたテキストを取得する
func getText(text string) string {
	m := strings.Split(strings.TrimSpace(text), " ")[1:]
	if len(m) == 0 {
		log.Fatal("invalid message")
	}
	return m[0]
}
