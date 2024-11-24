package main

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotLatestFeed = errors.New("not latest feed")
)

func CompareDiffFeed(ctx context.Context, latestFeed *AtomFeed) (*Entry, error) {
	// 今日の日付と比較して、今日以降のフィードを取得
	isLatest, err := isCompareDiffFeed(ctx, &latestFeed.Entries[0])
	if err != nil {
		return nil, err
	}

	if isLatest {
		return &latestFeed.Entries[0], nil
	}
	return nil, ErrNotLatestFeed
}

func isCompareDiffFeed(_ context.Context, feedEntity *Entry) (bool, error) {
	// ctxtime(https://github.com/newmo-oss/ctxtime)を使いたい
	// が、今回は、ctxtimeを使わずに、time.Timeを使う
	now := time.Now()
	if feedEntity.Published.Year() >= now.Year() && feedEntity.Published.Month() >= now.Month() && feedEntity.Published.Day() >= now.Day() {
		return true, nil
	}
	return false, nil
}
