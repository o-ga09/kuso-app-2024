package parser

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/newmo-oss/ctxtime/ctxtimetest"
	"github.com/newmo-oss/testid"
)

func TestCompareDiffFeed(t *testing.T) {
	type args struct {
		ctx        context.Context
		latestFeed *AtomFeed
		now        time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    *Entry
		wantErr error
	}{
		{name: "正常系 - 最新のフィードを取得できた場合", args: args{ctx: context.Background(), latestFeed: &AtomFeed{Entries: []Entry{{Published: time.Date(2024, 11, 11, 9, 0, 0, 0, time.UTC)}}}, now: time.Date(2024, 11, 11, 9, 0, 0, 0, time.UTC)}, want: &Entry{Published: time.Date(2024, 11, 11, 9, 0, 0, 0, time.UTC)}, wantErr: nil},
		{name: "正常系 - 最新のフィードを取得できなかった場合", args: args{ctx: context.Background(), latestFeed: &AtomFeed{Entries: []Entry{{Published: time.Date(2024, 11, 11, 9, 0, 0, 0, time.UTC)}}}, now: time.Date(2024, 11, 26, 9, 0, 0, 0, time.UTC)}, want: nil, wantErr: ErrNotLatestFeed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tid := uuid.NewString()
			ctx := testid.WithValue(tt.args.ctx, tid)
			ctxtimetest.SetFixedNow(t, ctx, tt.args.now)
			got, err := CompareDiffFeed(ctx, tt.args.latestFeed)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CompareDiffFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompareDiffFeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isCompareDiffFeed(t *testing.T) {
	type args struct {
		ctx        context.Context
		feedEntity *Entry
		now        time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "正常系 - 最新記事がツールが実行された当日だった場合、trueである", args: args{ctx: context.Background(), feedEntity: &Entry{Published: time.Date(2024, 11, 11, 9, 0, 0, 0, time.UTC)}, now: time.Date(2024, 11, 11, 9, 0, 0, 0, time.UTC)}, want: true},
		{name: "正常系 - 最新記事がツールが実行された当日以降だった場合、falseである", args: args{ctx: context.Background(), feedEntity: &Entry{Published: time.Date(2024, 11, 11, 9, 0, 0, 0, time.UTC)}, now: time.Date(2024, 11, 26, 9, 0, 0, 0, time.UTC)}, want: false},
		{name: "正常系 - 最新記事がツールが実行された日付の前日以降だった場合、trueである", args: args{ctx: context.Background(), feedEntity: &Entry{Published: time.Date(2024, 11, 10, 9, 0, 0, 0, time.UTC)}, now: time.Date(2024, 11, 11, 9, 0, 0, 0, time.UTC)}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tid := uuid.NewString()
			ctx := testid.WithValue(tt.args.ctx, tid)
			ctxtimetest.SetFixedNow(t, ctx, tt.args.now)
			got := isCompareDiffFeed(ctx, tt.args.feedEntity)
			if got != tt.want {
				t.Errorf("isCompareDiffFeed() = %v, want %v", got, tt.want)
			}
		})
	}
}
