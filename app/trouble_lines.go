package apps

import (
    "net/url"
    "github.com/PuerkitoBio/goquery"
)

type LineInfo struct {
    Name        string
    Outline     string
    Details     string
    Uri         string
}

type LineInfos []LineInfo

// Yahoo!路線情報から遅延している路線を取得
func getTroubleLines() LineInfos {
    _url := "https://transit.yahoo.co.jp/traininfo/area/4/" // 関東の路線情報

    doc, err := goquery.NewDocument(_url)
    if err != nil {
        panic(err)
    }

    u := url.URL{}
    u.Scheme = doc.Url.Scheme
    u.Host = doc.Url.Host

    lineInfos := LineInfos{};
    doc.Find("div.trouble > table > tbody").Each(func(_ int, s *goquery.Selection) {
        s.Children().Each(func(idx int, ss *goquery.Selection) {
            if (idx > 0) {
                lineInfo := LineInfo{}
                railway := ss.Children().Find("a").Text()
                href,_ := ss.Children().Find("a").Attr("href")
                status := ss.Children().Find("span.colTrouble").Text()
                detail := ss.Children().Next().Next().Text()
                // message = message + status + " @ " + railway + " (" + detail + ")\n"
                lineInfo.Name = railway
                lineInfo.Outline = status
                lineInfo.Details = detail
                lineInfo.Uri = href
                lineInfos = append(lineInfos, lineInfo)
            }
        })
    })
    return lineInfos
}
