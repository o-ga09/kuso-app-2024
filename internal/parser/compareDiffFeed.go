package parser

import (
	"context"
	"errors"

	"github.com/newmo-oss/ctxtime"
)

var (
	ErrNotLatestFeed = errors.New("not latest feed")
)

func CompareDiffFeed(ctx context.Context, latestFeed *AtomFeed) (*Entry, error) {
	// 今日の日付と比較して、今日以降のフィードを取得
	isLatest := isCompareDiffFeed(ctx, &latestFeed.Entries[0])
	if isLatest {
		return &latestFeed.Entries[0], nil
	}
	return nil, ErrNotLatestFeed
}

func isCompareDiffFeed(ctx context.Context, feedEntity *Entry) bool {
	now := ctxtime.Now(ctx)
	if feedEntity.Published.Year() >= now.Year() && feedEntity.Published.Month() >= now.Month() && feedEntity.Published.Day() >= now.Day() {
		return true
	}
	return false
}
