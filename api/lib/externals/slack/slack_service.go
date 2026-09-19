package slack

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

type SlackService struct {
	client    *slack.Client
	channelID string
}

// MessageParams BuildMessageBlocksに渡すメッセージパラメータ
type MessageParams struct {
	Title    string
	UserName string
	Date     string
	Weather  string
	Changes  string
}

// NewSlackService SlackServiceを初期化
func NewSlackService() (*SlackService, error) {
	// 環境変数が既に設定されている場合は.envファイルの読み込みをスキップ
	// docker-compose.ymlなどでenv_fileが設定されている場合は不要
	goEnv := os.Getenv("GO_ENV")
	if goEnv != "" {
		_ = godotenv.Load(fmt.Sprintf("../%s.env", goEnv))
		// エラーは無視（環境変数が既に設定されている場合があるため）
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("SLACK_BOT_TOKEN environment variable is not set")
	}

	// チャンネル送信は現在無効のため、SLACK_CHANNEL_IDのバリデーションは不要
	channelID := os.Getenv("SLACK_CHANNEL_ID")

	return &SlackService{
		client:    slack.New(botToken),
		channelID: channelID,
	}, nil
}

// SendMessage DMにメッセージを送信
func (s *SlackService) SendMessage(blocks []slack.Block, slackUserID string) error {
	// MTの議論により、チャンネルへのシフト変更通知は導入しない方針
	// チャンネル送信が必要になった場合はここを有効化する
	// _, _, err := s.client.PostMessage(
	// 	s.channelID,
	// 	slack.MsgOptionBlocks(blocks...),
	// )

	// 本人にDM送信 (IDがある場合のみ)
	if slackUserID != "" {
		_, _, err := s.client.PostMessage(
			slackUserID,
			slack.MsgOptionBlocks(blocks...),
		)
		if err != nil {
			var rateErr *slack.RateLimitedError
			if errors.As(err, &rateErr) {
				time.Sleep(rateErr.RetryAfter)
				// もう一度PostMessageを呼ぶ
				_, _, err = s.client.PostMessage(
					slackUserID,
					slack.MsgOptionBlocks(blocks...),
				)
				if err != nil {
					return fmt.Errorf("dm send error after retry: %w", err)
				}
			} else {
				return fmt.Errorf("dm send error: %w", err)
			}
		}
	}

	return nil
}

// BuildMessageBlocks リッチなメッセージを作成
func (s *SlackService) BuildMessageBlocks(params MessageParams) []slack.Block {
	headerText := fmt.Sprintf("🔔 %s", params.Title)
	headerBlock := slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", headerText, false, false))

	// 基本情報
	fields := []*slack.TextBlockObject{
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("ユーザー: %s", params.UserName), false, false),
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("日付: %s", params.Date), false, false),
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("天気: %s", params.Weather), false, false),
	}

	sectionBlock := slack.NewSectionBlock(nil, fields, nil)

	blocks := []slack.Block{headerBlock, sectionBlock}

	// 変更内容がある場合は追加
	if params.Changes != "" {
		changesBlock := slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*変更内容*\n%s", params.Changes), false, false),
			nil,
			nil,
		)
		blocks = append(blocks, changesBlock)
	}

	dividerBlock := slack.NewDividerBlock() // 区切り線
	blocks = append(blocks, dividerBlock)

	return blocks
}

// RescueMessageParams BuildRescueMessageBlocksに渡すメッセージパラメータ
type RescueMessageParams struct {
	Title     string   // 見出し(絵文字込み)
	Number    string   // 対応番号(T12 / Q3 / S5)。アプリの「本部からの返答」タブの表記に合わせる
	TypeLabel string   // トラブル / 質問 / 人が来ない
	Status    string   // 未対応 / 対応中 / 対応済
	Time      string   // 送信時刻(JST)
	Details   []string // 送信内容。「発生タスク: 受付」のような1行ずつ
	Response  string   // 本部からの返答(空なら載せない)
}

// mrkdwnで意味を持つ文字をエスケープする。送信内容や返答は利用者の入力なので、
// < や & がそのまま渡るとリンクやメンションとして解釈され表示が崩れる
var mrkdwnEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")

func escapeMrkdwn(s string) string {
	return mrkdwnEscaper.Replace(s)
}

// BuildRescueMessageBlocks レスキューの対応状況を知らせるメッセージを作成
func (s *SlackService) BuildRescueMessageBlocks(params RescueMessageParams) []slack.Block {
	headerBlock := slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", params.Title, false, false))

	fields := []*slack.TextBlockObject{
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("対応番号: %s", params.Number), false, false),
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("種類: %s", params.TypeLabel), false, false),
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("対応状況: %s", params.Status), false, false),
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("送信時刻: %s", params.Time), false, false),
	}
	blocks := []slack.Block{headerBlock, slack.NewSectionBlock(nil, fields, nil)}

	if len(params.Details) > 0 {
		lines := make([]string, len(params.Details))
		for i, d := range params.Details {
			lines[i] = escapeMrkdwn(d)
		}
		blocks = append(blocks, slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*送信内容*\n%s", strings.Join(lines, "\n")), false, false),
			nil,
			nil,
		))
	}

	if params.Response != "" {
		blocks = append(blocks, slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*本部からの返答*\n%s", escapeMrkdwn(params.Response)), false, false),
			nil,
			nil,
		))
	}

	blocks = append(blocks,
		slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", "アプリの「本部からの返答」タブでも確認できます", false, false)),
		slack.NewDividerBlock(),
	)
	return blocks
}

// SendRescueMessage レスキューの対応状況を送信者本人にDMする
func (s *SlackService) SendRescueMessage(params RescueMessageParams, slackUserID string) error {
	return s.SendMessage(s.BuildRescueMessageBlocks(params), slackUserID)
}
