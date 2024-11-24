package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

const FEED_URL = "https://go.dev/blog/feed.atom"

// AtomFeed はAtomフィードの構造体です
type AtomFeed struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Entries []Entry  `xml:"entry"`
}

// Entry はAtomフィードのエントリの構造体です
type Entry struct {
	Title     string    `xml:"title"`
	Link      Link      `xml:"link"`
	Published time.Time `xml:"published"`
	Summary   string    `xml:"summary"`
}

// Link はエントリのリンクの構造体です
type Link struct {
	Href string `xml:"href,attr"`
}

func ParseAtomFeed(ctx context.Context) (*AtomFeed, error) {
	// AtomフィードのURLを指定
	url := FEED_URL

	// フィードを取得
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching the feed:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// フィードを解析
	var feed AtomFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		fmt.Println("Error decoding the feed:", err)
		return nil, err
	}
	return &feed, nil
}
