package app

import (
	"github.com/PuerkitoBio/goquery"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
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

// 表示するテキスト
func getMessage() string {
	noTroubleMessage := "現在、遅延や運転の見合わせ等は発生していません。"

	// トラブルが発生している関東の路線を取得
	lineInfos := getTroubleLines()
	if lineInfos == nil {
		// トラブル無し
		return noTroubleMessage
	}

	// 表示対象の路線を取得
	targetLines, err := getTargetLines()
	if err != nil {
		panic(err)
	}

	var message string
	for _, tal := range lineInfos {
		if tal.isEnable(&targetLines) {
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
func (l *lineInfo) isEnable(t *targetLines) bool {
	if len(t.Lines) == 0 {
		// フィルターなし
		return true
	}
	for _, name := range t.Lines {
		if l.Name == name {
			return true
		}
	}
	return false
}

// Yahoo!路線情報から遅延している路線を取得
func getTroubleLines() []lineInfo {
	_url := "https://transit.yahoo.co.jp/traininfo/area/4/" // 関東の路線情報

	doc, err := goquery.NewDocument(_url)
	if err != nil {
		panic(err)
	}

	u := url.URL{}
	u.Scheme = doc.Url.Scheme
	u.Host = doc.Url.Host

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
