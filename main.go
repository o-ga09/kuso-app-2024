package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/o-ga09/kuso-app-2024/internal/genkit"
	"github.com/o-ga09/kuso-app-2024/internal/notifier"
	"github.com/o-ga09/kuso-app-2024/internal/parser"
)

const (
	RETRY_COUNT = 3
)

func main() {
	ctx := context.Background()

	// フィードを取得する
	feed, err := parser.ParseAtomFeed(ctx)
	if err != nil {
		fmt.Println("Error parsing the feed")
		os.Exit(1)
	}

	// 今日の日付と比較して、最新のフィードを取得する
	latestFeed, err := parser.CompareDiffFeed(ctx, feed)
	if err != nil && !errors.Is(err, parser.ErrNotLatestFeed) {
		log.Fatalf("Error comparing the feed: %v", err)
	}

	var message notifier.SlackMessage
	// ブログ記事をGeminiで要約する
	if errors.Is(err, parser.ErrNotLatestFeed) {
		// Slackに通知する
		message = notifier.SlackMessage{
			Text: "<!channel>\n 新しいブログ記事はありません",
		}
	} else {
		summary, err := genkit.SummarizeBlog(ctx, latestFeed.Link.Href)
		if err != nil {
			log.Fatalf("Error summarizing the blog: %v", err)
		}
		// Slackに通知する
		message = notifier.SlackMessage{
			Text: "<!channel>\n ⭐️ 新しいブログ記事が投稿されました！ ⭐️\nTitle: " + latestFeed.Title + "\nPublished: " + latestFeed.Published.Local().String() + "\n\n ======== \n" + summary + "========",
		}
	}

	// HTTPステータスが200以外の場合のみ3回リトライする
	for i := 0; i < RETRY_COUNT; i++ {
		err = notifier.SendSlackNotification(ctx, message)
		if err == nil || !errors.Is(err, notifier.ErrHTTPStatusNotOK) {
			break
		}
	}
	if err != nil {
		log.Fatalf("Error sending a Slack notification: %v", err)
	}
}
