package parser

import (
	"context"
	"encoding/xml"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestParseAtomFeed(t *testing.T) {
	type args struct {
		ctx        context.Context
		envVar     string
		mockResp   string
		mockStatus int
	}
	tests := []struct {
		name    string
		args    args
		want    *AtomFeed
		wantErr bool
	}{
		{
			name: "正常系 - Atomフィードを取得できた場合",
			args: args{
				ctx:    context.Background(),
				envVar: "https://example.com",
				mockResp: `<feed>
								<title>Example Feed</title>
								<entry>
									<title>Entry 1</title>
									<link href="http://example.com/entry1"/>
									<published>2023-10-01T12:00:00Z</published>
									<summary>Summary of entry 1</summary>
								</entry>
							</feed>`,
				mockStatus: http.StatusOK,
			},
			want: &AtomFeed{
				XMLName: xml.Name{Local: "feed"},
				Title:   "Example Feed",
				Entries: []Entry{
					{
						Title:     "Entry 1",
						Link:      Link{Href: "http://example.com/entry1"},
						Published: time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
						Summary:   "Summary of entry 1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "異常系 - 環境変数が空の場合",
			args: args{
				ctx:        context.Background(),
				envVar:     "",
				mockResp:   "",
				mockStatus: http.StatusOK,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系 - フィードの取得に失敗した場合",
			args: args{
				ctx:    context.Background(),
				envVar: "https://example.com",
				mockResp: `<feed>
								<title>Example Feed</title>
								<entry>
									<title>Entry 1</title>
									<link href="http://example.com/entry1"/>
									<published>2023-10-01T12:00:00Z</published>
									<summary>Summary of entry 1</summary>
								</entry>
							</feed>`,
				mockStatus: http.StatusInternalServerError,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系 - フィードのデコードに失敗した場合",
			args: args{
				ctx:        context.Background(),
				envVar:     "https://example.com",
				mockResp:   "",
				mockStatus: http.StatusOK,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("FEED_URL", tt.args.envVar)
			httpmock.Activate()
			t.Cleanup(httpmock.DeactivateAndReset)
			httpmock.RegisterResponder("GET", "https://example.com", httpmock.NewStringResponder(tt.args.mockStatus, tt.args.mockResp))
			got, err := ParseAtomFeed(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAtomFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAtomFeed() = %v, want %v", got, tt.want)
			}
		})
	}
}
