package app

import (
	"io/ioutil"
	"log"
	"regexp"
	"strconv"

	db "github.com/higuching/slack_bot/db"

	"github.com/PuerkitoBio/goquery"
	yaml "gopkg.in/yaml.v2"
)

type RailWays struct {
	data targetLinesInfo
}

type targetLinesInfo struct {
	Url    string   `yaml:"url"`
	Filter bool     `yaml:"filter"`
	Lines  []string `yaml:"lines"`
}

type lineInfo struct {
	Id      int
	Name    string
	Outline string
	Details string
	Url     string
}

// インスタンス化する
func NewRailWays() *RailWays {
	buf, err := ioutil.ReadFile("configs/railways.yml")
	if err != nil {
		log.Fatal(err)
	}

	o := RailWays{}
	err = yaml.Unmarshal(buf, &o.data)
	if err != nil {
		log.Fatal(err)
	}
	return &o
}

// 路線情報のテキストを取得する
func (o *RailWays) GetMessage() string {

	if o.data.Filter == false {
		// フィルター設定がOFFになっている
		//return "フィルター設定がOFFになっています"
		return ""
	}

	if len(o.data.Lines) == 0 {
		// 対象路線が設定されていない
		//return "通知対象の路線がOFFになっています"
		return ""
	}

	message := ""

	// DBインスタンス作成
	db := db.NewRailways()

	// トラブルが発生している関東の路線を取得
	troubleLines := getTroubleLines(o.data.Url)
	if troubleLines == nil {
		// トラブル無し
		rs := db.GetAll()
		if rs == nil {
			//return "現在、遅延や運転の見合わせ等は発生していません。"
			return ""
		}
		for _, r := range rs {
			// 登録済みの路線を解除する
			message = message + r.NAME + "の遅延が解消しました。\n"
		}
		// 全レコードの削除
		_ = db.DeleteAll()
		return message
	}

	// トラブっている路線情報を取得
	for _, tal := range troubleLines {
		if !tal.containsLine(&o.data) {
			// 対象路線に含まれる名前じゃない
			continue
		}
		if db.Get(tal.Id) {
			// レコードあるならすでに登録済み
			continue
		}
		_ = db.Insert(tal.Id, tal.Name)
		message = message + tal.Name + "で *" + tal.Outline + "* が発生しました。 " + tal.Url + "" + "\n"
	}

	// トラブルが解消した路線情報を取得
	rs := db.GetAll()
	if rs != nil {
		for _, r := range rs {
			isFind := false
			for _, tal := range troubleLines {
				if r.ID == tal.Id {
					isFind = true
				}
			}
			if isFind == false {
				// 解消！
				_ = db.Delete(r.ID)
				message = message + r.NAME + "の遅延が解消しました。\n"
			}
		}
	}

	if message == "" {
		// 指定の路線でトラブル無し
		//return "新規に遅延や運転の見合わせ等は発生していませんでした。"
		return ""
	}
	return message
}

// 必要な路線か判定
func (l *lineInfo) containsLine(t *targetLinesInfo) bool {
	for _, name := range t.Lines {
		if l.Name == name {
			return true
		}
	}
	return false
}

// Yahoo!路線情報から遅延している路線を取得
func getTroubleLines(_url string) []lineInfo {
	doc, err := goquery.NewDocument(_url)
	if err != nil {
		panic(err)
	}

	li := []lineInfo{}
	doc.Find("div.trouble > table > tbody").Each(func(_ int, s *goquery.Selection) {
		s.Children().Each(func(idx int, ss *goquery.Selection) {
			if idx > 0 {
				href, _ := ss.Children().Find("a").Attr("href")
				// URLからIDを抽出
				r := regexp.MustCompile(`[\d]+`)
				slice := r.FindAllStringSubmatch(href, -1)
				id, err2 := strconv.Atoi(slice[0][0])
				if err2 != nil {
					panic(err2)
				}
				li = append(li, lineInfo{
					Id:      id,
					Name:    ss.Children().Find("a").Text(),
					Outline: ss.Children().Find("span.colTrouble").Text(),
					Details: ss.Children().Next().Next().Text(),
					Url:     href,
				})
			}
		})
	})
	return li
}
