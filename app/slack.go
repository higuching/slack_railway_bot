package apps

import (
    "log"
    "os"
    "encoding/json"
    "io/ioutil"
    "github.com/nlopes/slack"
)

type SlackJson struct {
    Token   string
}

// func Run(api *slack.Client) int {
func Run() {
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
                log.Print("Hello Event")

            case *slack.MessageEvent:
                switch ev.Text {
                case "運行状況":
                    message := GetMessage()
                    rtm.SendMessage(rtm.NewOutgoingMessage(message, ev.Channel))
                case "help":
                    rtm.SendMessage(rtm.NewOutgoingMessage("Usage: 運行状況　関東圏の電車の運行状況を表示します。\n他の機能？ないよ。", ev.Channel))
                default:
                    // rtm.SendMessage(rtm.NewOutgoingMessage("登録されていないパターンです", ev.Channel))
                }

            case *slack.InvalidAuthEvent:
                log.Print("Invalid credentials")
                os.Exit(1)

            }
        }
    }
}

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
