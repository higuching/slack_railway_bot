package app

import (
	"github.com/PuerkitoBio/goquery"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type targetLines struct {
	Filter bool     `yaml:"filter"`
	Lines  []string `yaml:"lines"`
}

type lineInfo struct {
	Name    string
	Outline string
	Details string
	Uri     string
}

const noTroubleMessage = "現在、遅延や運転の見合わせ等は発生していません。"

// 表示するテキスト
func getMessage() string {

	// トラブルが発生している関東の路線を取得
	_url := "https://transit.yahoo.co.jp/traininfo/area/4/" // 関東の路線情報
	troubleLines := getTroubleLines(_url)
	if troubleLines == nil {
		// トラブル無し
		return noTroubleMessage
	}

	// 表示対象の路線を取得
	targetLines, err := getTargetLines()
	if err != nil {
		panic(err)
	}

	var message string
	for _, tal := range troubleLines {
		if len(targetLines.Lines) == 0 || tal.containsLine(&targetLines) {
			// フィルターなし or 対象路線に含まれる名前
			message = message + tal.Name + " @ " + tal.Outline + "(" + tal.Details + ")" + "\n"
		}
	}
	if message == "" {
		// 指定の路線でトラブル無し
		return noTroubleMessage
	}
	return message
}

// 定義した情報の欲しい路線情報を取得
func getTargetLines() (targetLines, error) {
	buf, err := ioutil.ReadFile("configs/railways.yml")
	if err != nil {
		return targetLines{}, err
	}

	tl := targetLines{}
	err = yaml.Unmarshal(buf, &tl)
	if err != nil {
		return targetLines{}, err
	}
	return tl, nil
}

// 必要な路線か判定
func (l *lineInfo) containsLine(t *targetLines) bool {
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
				li = append(li, lineInfo{
					Name:    ss.Children().Find("a").Text(),
					Outline: ss.Children().Find("span.colTrouble").Text(),
					Details: ss.Children().Next().Next().Text(),
					Uri:     href,
				})
			}
		})
	})
	return li
}
