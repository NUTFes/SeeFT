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
	Title     string
	UserName  string
	Date      string
	Weather   string
	TimeRange string
	Changes   string
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

	channelID := os.Getenv("SLACK_CHANNEL_ID")
	if channelID == "" {
		return nil, fmt.Errorf("SLACK_CHANNEL_ID environment variable is not set")
	}

	return &SlackService{
		client:    slack.New(botToken),
		channelID: channelID,
	}, nil
}

// SendMessage チャンネルとDMにメッセージを送信
func (s *SlackService) SendMessage(blocks []slack.Block, slackUserID string) error {
	// 1. チャンネルに送信
	_, _, err := s.client.PostMessage(
		s.channelID,
		slack.MsgOptionBlocks(blocks...),
	)
	if err != nil {
		var rateErr *slack.RateLimitedError
		if errors.As(err, &rateErr) {
			time.Sleep(rateErr.RetryAfter)
			// もう一度PostMessageを呼ぶ
			_, _, err = s.client.PostMessage(
				s.channelID,
				slack.MsgOptionBlocks(blocks...),
			)
			if err != nil {
				return fmt.Errorf("channel send error after retry: %w", err)
			}
		} else {
			return fmt.Errorf("channel send error: %w", err)
		}
	}

	// 2. 本人にDM送信 (IDがある場合のみ)
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
					return fmt.Errorf("channel send error after retry: %w", err)
				}
			} else {
				return fmt.Errorf("channel send error: %w", err)
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

	// 時間範囲がある場合は追加
	if params.TimeRange != "" {
		fields = append(fields,
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("時間: %s", params.TimeRange), false, false),
		)
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
