package genkit

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func Test_fetchWebContent(t *testing.T) {
	type args struct {
		url        string
		mockStatus int
		mockResp   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "正常系　- Webページのコンテンツを取得できた場合",
			args: args{
				url:        "https://example.com",
				mockStatus: http.StatusOK,
				mockResp: `<html>
								<head>
									<title>Google</title>
								</head>	
								<body>
									<h1>Google</h1>
								</body>
							</html>`,
			},
			want:    "Google",
			wantErr: false,
		},
		{
			name: "異常系 - URLが不正な場合",
			args: args{
				url:        "https://example.com",
				mockStatus: http.StatusNotFound,
				mockResp:   "",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "正常系 - script, style, noscriptタグが削除されていること",
			args: args{
				url:        "https://example.com",
				mockStatus: http.StatusOK,
				mockResp: `<html>
								<body>
									<h1>Title</h1>
									<script>alert('Hello, World!');</script>
									<style>h1 { color: red; }</style>
									<noscript><h1>JavaScript is disabled</h1></noscript>
								</body>
							</html>`,
			},
			want:    "Title",
			wantErr: false,
		},
		{
			name: "正常系 - articleタグが存在する場合",
			args: args{
				url:        "https://example.com",
				mockStatus: http.StatusOK,
				mockResp: `<html>
								<body>	
									<article>
										<h1>Article Title</h1>
									</article>
								</body>
							</html>`,
			},
			want:    "Article Title",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			t.Cleanup(httpmock.DeactivateAndReset)
			httpmock.RegisterResponder("GET", "https://example.com", httpmock.NewStringResponder(tt.args.mockStatus, tt.args.mockResp))
			got, err := fetchWebContent(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchWebContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fetchWebContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
