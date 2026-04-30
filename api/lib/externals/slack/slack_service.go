package slack

import (
	"errors"
	"fmt"
	"os"
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
