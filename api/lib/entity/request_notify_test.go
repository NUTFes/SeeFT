package entity

import (
	"encoding/json"
	"testing"
)

// 本部の書き込み(onChange)は notify を送らない。省略を「知らせない」と読むと、
// 通常の対応状況・返答の変更でDMが一切届かなくなる
func TestRescueUpdateRequest_notifyを省略したら知らせる(t *testing.T) {
	cases := []struct {
		name string
		body string
		want bool
	}{
		{"本部の書き込み(notifyなし)", `{"status":"inProgress","response":""}`, true},
		{"notify: true", `{"status":"done","response":"","notify":true}`, true},
		{"GASが重複をまとめる書き込み(notify: false)", `{"status":"done","response":"対応番号6にまとめました","notify":false}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var req TroubleRescueUpdateRequest
			if err := json.Unmarshal([]byte(tc.body), &req); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got := req.ShouldNotify(); got != tc.want {
				t.Errorf("ShouldNotify() = %v, want %v", got, tc.want)
			}
		})
	}
}
