package usecase

import (
	"testing"

	"github.com/NUTFes/SeeFT/api/lib/entity" 
)

func Test_notificationUseCase_formatTimeRange(t *testing.T) {
	type args struct {
		timeMap      map[int]entity.Time
		startTimeID  int
		endTimeID    int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "正常系: マップに存在するIDが正しくフォーマットされること",
			args: args{
				timeMap: map[int]entity.Time{
					1: {Time: "10:00"},
					3: {Time: "12:00"}, // endTimeID(2) + 1 = 3 の場所
				},
				startTimeID: 1,
				endTimeID:   2,
			},
			want: "10:00 〜 12:00",
		},
		{
			name: "正常系: 終了ID+1がマップに存在しない場合、0:00にフォールバックすること",
			args: args{
				timeMap: map[int]entity.Time{
					1: {Time: "23:00"},
					// ID: 3 (2+1) が存在しない
				},
				startTimeID: 1,
				endTimeID:   2,
			},
			want: "23:00 〜 0:00",
		},
		{
			name: "要判断: 開始IDがマップに存在しない場合、開始時刻が空文字になる（非対称な挙動: issue #418）",
			args: args{
				timeMap: map[int]entity.Time{
					3: {Time: "12:00"},
					// ID: 1 が存在しない
				},
				startTimeID: 1,
				endTimeID:   2,
			},
			want: " 〜 12:00", // 開始側が空文字のため、先頭にスペースが残る現状の挙動
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// レシーバはフィールドを参照しない純関数なのでゼロ値で初期化
			n := &notificationUseCase{}
			
			got := n.formatTimeRange(tt.args.timeMap, tt.args.startTimeID, tt.args.endTimeID)
			if got != tt.want {
				t.Errorf("formatTimeRange() = %v, want %v", got, tt.want)
			}
		})
	}
}