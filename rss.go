package main

import (
    "github.com/mmcdole/gofeed"
)

// RSSから取得するタイプのサンプル
func getRailwayRss() string {
    fp := gofeed.NewParser()

    uri := "http://api.tetsudo.com/traffic/rss20.xml?kanto";
    feed, _ := fp.ParseURL(uri)
    items := feed.Items

    if (items == nil) {
        return ""
    }

    var message string
    for _, item := range items {
        message = message + item.Title + item.Description + "\n"
    }
    return message
}
