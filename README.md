# クソアプリハッカソン2024

[クソアプリアドベントカレンダー](https://qiita.com/advent-calendar/2024/kuso-app)のアプリをハッカソン形式で作る！2024


[クソアプリハッカソン 2024 @ Findy](https://kuso-app.connpass.com/event/336557/)

[Qiita記事(TBD)](https://qiita.com/o-ga)

[発表資料](https://docs.google.com/presentation/d/12SoPi9srlx3E_gwxCvr6YAPmIsZPNJK8nYJ-H7tyOjw/edit?usp=sharing)

---

## つくるもの

「Go BlogをLLMで要約して、翻訳して、通知する」

## 技術構成

- [Go](https://go.dev/)
- [Gemini](https://gemini.google.com/)
- [Genkit](https://firebase.google.com/products/genkit?hl=ja)

## CI/CD

- GitHub Actionsを使用して定期的に実行

## 環境変数

- GOOGLE_GENAI_API_KEY: GoogleのLLM APIキー
- SLACK_WEBHOOK_URL: SlackのWebhook URL

## 実行方法
- 実行環境に以下の環境変数を設定します。
- 以下のコマンドを実行してアプリケーションを起動します。
  
```bash
$ go run .
```

## GitHub Actions

GitHub Actionsを使用して、毎日午前8時30分(GitHub Actionsは無料枠なので、午前9時ごろに起動する)にアプリケーションを実行します。設定は cron.yml に記載されています。

## Author

@o-ga09
