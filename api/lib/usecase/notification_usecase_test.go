package usecase

import (
	"testing"

	"github.com/NUTFes/SeeFT/api/lib/entity"
)

func Test_notificationUseCase_formatTimeRange(t *testing.T) {
	// テスト対象の構造体をゼロ値で生成
	n := &notificationUseCase{}

	tests := []struct {
		name        string
		timeMap     map[int]entity.Time
		startTimeID int
		endTimeID   int
		want        string
	}{
		{
			name: "正常系: 連続スロット範囲のフォーマット",
			timeMap: map[int]entity.Time{
				1: {Time: "9:00"},
				2: {Time: "10:00"},
				3: {Time: "11:00"},
			},
			startTimeID: 1,
			endTimeID:   2,
			want:        "9:00 〜 11:00",
		},
		{
			name: "正常系: 同一スロット指定（実際の呼び出し形態）",
			timeMap: map[int]entity.Time{
				5: {Time: "13:00"},
				6: {Time: "14:00"},
			},
			startTimeID: 5,
			endTimeID:   5,
			want:        "13:00 〜 14:00",
		},
		{
			name:        "境界値: 最終スロットで次スロットが無い（要素1個のマップ）",
			timeMap:     map[int]entity.Time{10: {Time: "22:00"}},
			startTimeID: 10,
			endTimeID:   10,
			want:        "22:00 〜 0:00",
		},
		{
			name:        "境界値: nil マップとゼロ値 ID",
			timeMap:     nil,
			startTimeID: 0,
			endTimeID:   0,
			want:        " 〜 0:00",
		},
		{
	
			name:        "要判断: startTimeID がマップに存在しない（非対称な挙動: issue #）",
			timeMap:     map[int]entity.Time{2: {Time: "10:00"}},
			startTimeID: 5,
			endTimeID:   1,
			want:        " 〜 10:00",
		},
		{
			name: "境界値: 次スロットは存在するが Time フィールドがゼロ値（空文字）",
			timeMap: map[int]entity.Time{
				1: {Time: "9:00"},
				2: {},
			},
			startTimeID: 1,
			endTimeID:   1,
			want:        "9:00 〜 ",
		},
		{
			name:        "異常系: 負の ID（endTimeID+1 の算術で key 0 を参照）",
			timeMap:     map[int]entity.Time{0: {Time: "8:00"}},
			startTimeID: -1,
			endTimeID:   -1,
			want:        " 〜 8:00",
		},
		{
			name: "異常系: 逆転範囲（startTimeID > endTimeID）",
			timeMap: map[int]entity.Time{
				1: {Time: "9:00"},
				2: {Time: "10:00"},
				3: {Time: "11:00"},
				4: {Time: "12:00"},
			},
			startTimeID: 3,
			endTimeID:   1,
			want:        "11:00 〜 10:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := n.formatTimeRange(tt.timeMap, tt.startTimeID, tt.endTimeID)
			if got != tt.want {
				t.Errorf("formatTimeRange() = %v, want %v", got, tt.want)
			}
		})
	}
}