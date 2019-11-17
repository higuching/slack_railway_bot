package app

import (
	"io/ioutil"
	"regexp"
	"strconv"

	// トラブルが発生している関東の路線情報を格納するDB
	db "github.com/higuching/slack_railway_bot/db"

	"github.com/PuerkitoBio/goquery"
	yaml "gopkg.in/yaml.v2"
)

// railwaysConfig 路線情報を取得する設定
type railwaysConfig struct {
	URL    string   `yaml:"url"`
	Filter bool     `yaml:"filter"`
	Lines  []string `yaml:"lines"`
}

// lineInfo 取得した路線情報
type lineInfo struct {
	ID      int
	Name    string
	Outline string
	Details string
	URL     string
}

// getRailwaysConfig 路線情報の設定ファイルを取得する
func getRailwaysConfig() (*railwaysConfig, error) {
	o := railwaysConfig{}

	buf, err := ioutil.ReadFile("configs/railways.yml")
	if err != nil {
		return &o, err
	}

	err = yaml.Unmarshal(buf, &o)
	if err != nil {
		return &o, err
	}
	return &o, nil
}

// getMessage HTMLをスクレイピングして路線情報のテキストを取得する
func getMessage() string {

	msg := ""
	o, err := getRailwaysConfig()
	if err != nil {
		panic(err)
	}

	if o.Filter == false {
		// フィルター設定がOFFになっている
		return msg
	}

	if len(o.Lines) == 0 {
		// 対象路線が設定されていない
		return msg
	}

	// DBインスタンス作成
	db := db.NewRailways()

	// トラブルが発生している関東の路線を取得
	troubleLines := getTroubleLines(o.URL)
	if troubleLines == nil {
		// トラブル無し
		rs := db.GetAll()
		if rs == nil {
			return ""
		}
		for _, r := range rs {
			// 登録済みの路線を解除する
			msg = msg + r.NAME + "の遅延が解消しました。\n"
		}
		// 全レコードの削除
		_ = db.DeleteAll()
		return msg
	}

	// トラブっている路線情報を取得
	for _, tal := range troubleLines {
		if !tal.containsLine(o) {
			// 対象路線に含まれる名前じゃない
			continue
		}
		if db.Get(tal.ID) {
			// レコードあるならすでに登録済み
			continue
		}
		_ = db.Insert(tal.ID, tal.Name)
		msg = msg + tal.Name + "で *" + tal.Outline + "* が発生しました。 " + tal.URL + "" + "\n"
	}

	// トラブルが解消した路線情報を取得
	rs := db.GetAll()
	if rs != nil {
		for _, r := range rs {
			isFind := false
			for _, tal := range troubleLines {
				if r.ID == tal.ID {
					isFind = true
				}
			}
			if isFind == false {
				// 解消！
				_ = db.Delete(r.ID)
				msg = msg + r.NAME + "の遅延が解消しました。\n"
			}
		}
	}

	if msg == "" {
		// 指定の路線でトラブル無し
		return ""
	}
	return msg
}

// containsLine 必要な路線か判定
func (l *lineInfo) containsLine(t *railwaysConfig) bool {
	for _, name := range t.Lines {
		if l.Name == name {
			return true
		}
	}
	return false
}

// getTroubleLines 遅延している路線を取得
func getTroubleLines(u string) []lineInfo {
	doc, err := goquery.NewDocument(u)
	if err != nil {
		panic(err)
	}

	li := []lineInfo{}
	doc.Find("div.trouble > table > tbody").Each(func(_ int, s *goquery.Selection) {
		s.Children().Each(func(idx int, ss *goquery.Selection) {
			if idx == 0 {
				return
			}
			href, _ := ss.Children().Find("a").Attr("href")
			// URLからIDを抽出
			r := regexp.MustCompile(`[\d]+`)
			slice := r.FindAllStringSubmatch(href, -1)
			id, err2 := strconv.Atoi(slice[0][0])
			if err2 != nil {
				panic(err2)
			}
			li = append(li, lineInfo{
				ID:      id,
				Name:    ss.Children().Find("a").Text(),
				Outline: ss.Children().Find("span.colTrouble").Text(),
				Details: ss.Children().Next().Next().Text(),
				URL:     href,
			})
		})
	})
	return li
}
