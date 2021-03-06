package app

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	db "github.com/higuching/slack_railway_bot/db"
	yaml "gopkg.in/yaml.v2"
)

// slackの設定ファイル
type slackConfig struct {
	Url     string `yaml:"url"`
	Channel string `yaml:"channel"`
	Nmae    string `yaml:"name"`
}

// mainから呼ばれるコアとなる関数
func Run() (string, error) {
	// DBの準備
	err := db.NewRailways().Create()
	if err != nil {
		return "", err
	}
	log.Print("Database is ready.")

	// slackの設定読み込み
	conf, err := getSlackConfig()
	if err != nil {
		return "", err
	}

	// POSTするテキストの取得
	text := getPostText(conf.Nmae, conf.Channel)
	if text == "" {
		// なければ終わり
		return "", err
	}

	// slackへPOST
	req, err := http.NewRequest(
		"POST",
		conf.Url,
		bytes.NewBuffer([]byte(text)),
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

// slackの設定ファイルを読み込む
func getSlackConfig() (*slackConfig, error) {
	o := slackConfig{}
	buf, err := ioutil.ReadFile("configs/slack.yml")
	if err != nil {
		return &o, err
	}

	err = yaml.Unmarshal(buf, &o)
	if err != nil {
		return &o, err
	}
	return &o, nil
}

// slackに投げるテキストを生成する
func getPostText(n, c string) string {
	t := getMessage()
	if t == "" {
		log.Print("Infomartion is not updated.")
		return ""
	}
	log.Printf("Infomartion is updated. Text:%s", t)

	return `{"channel":"` + c + `","username":"` + n + `","text":"` + t + `"}`
}
