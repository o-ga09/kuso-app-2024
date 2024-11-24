package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	RETRY_COUNT = 3
)

func main() {
	ctx := context.Background()

	// フィードを取得する
	feed, err := ParseAtomFeed(ctx)
	if err != nil {
		fmt.Println("Error parsing the feed")
		os.Exit(1)
	}

	// 今日の日付と比較して、最新のフィードを取得する
	latestFeed, err := CompareDiffFeed(ctx, feed)
	if err != nil && !errors.Is(err, ErrNotLatestFeed) {
		log.Fatalf("Error comparing the feed: %v", err)
	}

	var message SlackMessage
	// ブログ記事をGeminiで要約する
	if errors.Is(err, ErrNotLatestFeed) {
		// Slackに通知する
		message = SlackMessage{
			Text: "<!channel>\n 新しいブログ記事はありません",
		}
	} else {
		summary, err := SummarizeBlog(ctx, latestFeed.Link.Href)
		if err != nil {
			log.Fatalf("Error summarizing the blog: %v", err)
		}
		// Slackに通知する
		message = SlackMessage{
			Text: "<!channel>\n ⭐️ 新しいブログ記事が投稿されました！ ⭐️\n\n ======== \n" + summary + "========",
		}
	}

	// HTTPステータスが200以外の場合のみ3回リトライする
	for i := 0; i < RETRY_COUNT; i++ {
		err = SendSlackNotification(ctx, message)
		if err == nil || !errors.Is(err, ErrHTTPStatusNotOK) {
			break
		}
	}
	if err != nil {
		log.Fatalf("Error sending a Slack notification: %v", err)
	}
}
