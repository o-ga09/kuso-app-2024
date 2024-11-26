package parser

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"time"
)

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
	url := os.Getenv("FEED_URL")
	if url == "" {
		return nil, fmt.Errorf("FEED_URL is empty")
	}

	// フィードを取得
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching the feed: %v", err)
	}
	defer resp.Body.Close()

	// フィードを解析
	var feed AtomFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {

		return nil, fmt.Errorf("error decoding the feed: %v", err)
	}
	return &feed, nil
}
