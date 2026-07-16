package usecase

import "testing"

// compareTimeStrings のゴールデンテスト。
// 期待値は docs/development/test-design/phase1-pure-functions.md の
// compareTimeStrings セクションに記載のケース表と、実行による裏取り結果に基づく。
// 「要判断」の付いたケースは現状の挙動をそのまま固定したもので、
// 挙動の是非は issue #418 で判断する。
func TestCompareTimeStrings(t *testing.T) {
	u := &shiftUseCase{}

	cases := []struct {
		name  string
		time1 string
		time2 string
		want  int
	}{
		{
			name:  "桁数が異なる時刻の数値比較",
			time1: "8:00", time2: "10:00",
			want: -1,
		},
		{
			name:  "逆順で正の値を返す",
			time1: "10:00", time2: "8:00",
			want: 1,
		},
		{
			name:  "同一時刻は等価",
			time1: "9:15", time2: "9:15",
			want: 0,
		},
		{
			name:  "同時間帯での分単位の差",
			time1: "9:15", time2: "9:30",
			want: -1,
		},
		{
			name:  "ゼロ値時刻と一日の最大時刻",
			time1: "0:00", time2: "23:59",
			want: -1,
		},
		{
			name:  "ゼロパディング表記の同値性",
			time1: "08:05", time2: "8:05",
			want: 0,
		},
		{
			name:  "非正規化の分表記は換算後に等価",
			time1: "1:30", time2: "0:90",
			want: 0,
		},
		{
			// 要判断（issue #418）: 空文字列はガード節に入り「等価」扱いになる。
			// sort.Slice は非安定ソートのため、異常な StartTime を持つ
			// ShiftCard の並び順が実行のたびに変わりうる。
			name:  "空文字列は正当な時刻とも等価扱い（要判断: issue #418）",
			time1: "", time2: "10:00",
			want: 0,
		},
		{
			// 要判断（issue #418）: コロンが2個以上の入力も、空文字列と同じ
			// ガード節（len(parts) != 2）を通り「等価」扱いになる。上の
			// 「空文字列」ケースと同一の論点（issue #418 の論点1）。
			name:  "コロンが2個以上の書式は等価扱い（要判断: issue #418）",
			time1: "10:00:00", time2: "9:00",
			want: 0,
		},
		{
			// 要判断（issue #418）: strconv.Atoi のエラーが黙殺され、
			// 非数値の時分が 0:00 として扱われる。上の「空文字列」ケースとは
			// 異なり「等価」ではなく「最小値」として扱われる非対称な挙動。
			name:  "非数値の時分は暗黙に0:00扱い（要判断: issue #418）",
			time1: "aa:bb", time2: "8:00",
			want: -1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := u.compareTimeStrings(c.time1, c.time2)
			if got != c.want {
				t.Errorf("compareTimeStrings(%q, %q) = %d, want %d", c.time1, c.time2, got, c.want)
			}
		})
	}
}
