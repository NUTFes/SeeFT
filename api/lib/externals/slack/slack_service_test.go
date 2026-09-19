package slack

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildRescueMessageBlocks_利用者の入力をエスケープし返答が無ければ載せない(t *testing.T) {
	s := &SlackService{}

	withResponse := s.BuildRescueMessageBlocks(RescueMessageParams{
		Title:     "💬 本部から返答が届きました",
		Number:    "Q4",
		TypeLabel: "質問",
		Status:    "対応中",
		Time:      "2026/07/24 18:00:00",
		Details:   []string{"質問: A&B <どっち>?"},
		Response:  "<!channel> 体育館です",
	})
	b, err := json.Marshal(withResponse)
	require.NoError(t, err)
	body := string(b)
	assert.Contains(t, body, "A\\u0026amp;B \\u0026lt;どっち\\u0026gt;?")
	// <!channel> がそのまま渡るとチャンネル全員へのメンションとして解釈される
	assert.Contains(t, body, "\\u0026lt;!channel\\u0026gt; 体育館です")
	assert.Contains(t, body, "*本部からの返答*")

	withoutResponse := s.BuildRescueMessageBlocks(RescueMessageParams{Title: "👀 本部がレスキューを確認しました", Number: "T7"})
	b, err = json.Marshal(withoutResponse)
	require.NoError(t, err)
	assert.False(t, strings.Contains(string(b), "*本部からの返答*"))
	assert.Len(t, withoutResponse, len(withResponse)-2) // 送信内容と返答のセクションが無い
}
