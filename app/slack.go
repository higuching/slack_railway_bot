package apps

import (
    "log"
    "os"
    "encoding/json"
    "io/ioutil"
    "strings"
    "github.com/nlopes/slack"
)

type SlackJson struct {
    Token   string
}

var botId string

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

            case *slack.ConnectedEvent:
                botId = ev.Info.User.ID

            case *slack.MessageEvent:
                if strings.HasPrefix(ev.Text, "<@" + botId + ">") {
                    text := getText(ev.Text)
                    switch text {
                    case "運行状況":
                        message := GetMessage()
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
