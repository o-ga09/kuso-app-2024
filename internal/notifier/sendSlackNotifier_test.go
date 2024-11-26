package notifier

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestSendSlackNotification(t *testing.T) {
	type args struct {
		ctx     context.Context
		message SlackMessage
		envVar  string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{name: "正常系 - Slack通知できた場合", args: args{ctx: context.Background(), message: SlackMessage{Text: "Hello, World!"}, envVar: "http://example.com"}, wantErr: false, wantStatus: http.StatusOK},
		{name: "準正常系 - リクエストパラメータが異なる場合", args: args{ctx: context.Background(), message: SlackMessage{Text: "Hello, World!"}, envVar: "http://example.com"}, wantErr: true, wantStatus: http.StatusBadRequest},
		{name: "異常系 - 環境変数SLACK_WEBHOOK_URLが空文字の場合", args: args{ctx: context.Background(), message: SlackMessage{Text: "Hello, World!"}, envVar: ""}, wantErr: true, wantStatus: http.StatusBadRequest},
		{name: "異常系 - Slack通知するメッセージが空文字場合", args: args{ctx: context.Background(), message: SlackMessage{Text: ""}, envVar: "http://example.com"}, wantErr: false, wantStatus: http.StatusOK},
		{name: "異常系 - SlackのWebhookが500エラーを返した場合", args: args{ctx: context.Background(), message: SlackMessage{Text: "Hello, World!"}, envVar: "http://example.com"}, wantErr: true, wantStatus: http.StatusInternalServerError},
	}

	originalHTTPClient := http.DefaultClient
	defer func() { http.DefaultClient = originalHTTPClient }()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			t.Cleanup(httpmock.DeactivateAndReset)

			httpmock.RegisterResponder("POST", "http://example.com", httpmock.NewStringResponder(tt.wantStatus, ""))
			os.Setenv("SLACK_WEBHOOK_URL", tt.args.envVar)
			err := SendSlackNotification(tt.args.ctx, tt.args.message)
			t.Log(err)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendSlackNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
