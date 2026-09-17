package entity

import (
	"testing"
	"time"
)

// スプレッドシートへ送る発生時刻はコントローラがAsia/Tokyoで整形している。
// アプリへ返す時刻がUTCのままだと、同じレコードが9時間ずれて見える（#535）。
func TestFormatRescueTime_ReturnsJST(t *testing.T) {
	// 2026-09-17 15:16:56 UTC は JST では翌日の 00:16:56
	stored := time.Date(2026, 9, 17, 15, 16, 56, 0, time.UTC)

	got := formatRescueTime(stored)

	const want = "2026/09/18 00:16:56"
	if got != want {
		t.Fatalf("formatRescueTime(%s) = %q, want %q", stored, got, want)
	}
}

// DBから返る値のロケーションに関係なく、表示は常にJSTで揃える。
func TestFormatRescueTime_NormalizesOtherLocations(t *testing.T) {
	utc := time.Date(2026, 9, 17, 15, 16, 56, 0, time.UTC)
	sameInstantElsewhere := utc.In(time.FixedZone("UTC+7", 7*60*60))

	if got, want := formatRescueTime(sameInstantElsewhere), formatRescueTime(utc); got != want {
		t.Fatalf("同じ瞬間なのに表示が違う: %q と %q", got, want)
	}
}
