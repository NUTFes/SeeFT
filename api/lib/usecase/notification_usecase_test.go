package usecase

import (
	"reflect"
	"sort"
	"testing"

	"github.com/NUTFes/SeeFT/api/lib/entity"
)

// sortLogsByTime のゴールデンテスト。
// 期待値は docs/development/test-design/phase1-pure-functions.md の
// sortLogsByTime セクションに記載のケース表と、実行による裏取り結果に基づく。
// 「要判断」の付いたケースは現状の挙動をそのまま固定したもので、
// 挙動の是非は issue #TODO-A（サイレント除外）と issue #TODO-B（タイの順序）で判断する。
func TestSortLogsByTime(t *testing.T) {
	n := &notificationUseCase{}

	cases := []struct {
		name     string
		logs     []entity.ActionLog
		shiftMap map[int]entity.ShiftAdmin
		wantIDs  []int
		// sort.Slice は非安定ソートのため、同一 TimeID のケースだけは
		// 順序を問わず集合として比較する（issue #TODO-B）。
		unordered bool
	}{
		{
			name: "複数ログをTimeID昇順に並べ替える",
			logs: []entity.ActionLog{
				{ID: 1, ShiftID: 10},
				{ID: 2, ShiftID: 20},
				{ID: 3, ShiftID: 30},
			},
			shiftMap: map[int]entity.ShiftAdmin{
				10: {ID: 10, TimeID: 3},
				20: {ID: 20, TimeID: 1},
				30: {ID: 30, TimeID: 2},
			},
			wantIDs: []int{2, 3, 1},
		},
		{
			name:     "要素1個はそのまま返る",
			logs:     []entity.ActionLog{{ID: 1, ShiftID: 10}},
			shiftMap: map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 5}},
			wantIDs:  []int{1},
		},
		{
			name:     "空スライスは空の非nilスライスを返す",
			logs:     []entity.ActionLog{},
			shiftMap: map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 1}},
			wantIDs:  []int{},
		},
		{
			// logs はゼロ値（nil スライス）のまま渡す。
			name:     "nilスライスでもpanicせず空を返す",
			shiftMap: map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 1}},
			wantIDs:  []int{},
		},
		{
			// shiftMap はゼロ値（nil マップ）のまま渡す。nil マップの
			// lookup は ok=false になり、全件 continue される。
			name: "nilマップでは全ログが除外され空を返す",
			logs: []entity.ActionLog{
				{ID: 1, ShiftID: 10},
				{ID: 2, ShiftID: 20},
			},
			wantIDs: []int{},
		},
		{
			// 要判断（issue #TODO-A）: shiftMap にキーが無いログは、エラーにも
			// ログにもならず結果から黙って除外される。関数名は sort だが
			// フィルタを兼ねている現状挙動をそのまま固定したもの。
			name: "shiftMapにキーが無いログは黙って除外される（要判断: issue #TODO-A）",
			logs: []entity.ActionLog{
				{ID: 1, ShiftID: 10},
				{ID: 2, ShiftID: 99},
				{ID: 3, ShiftID: 30},
			},
			shiftMap: map[int]entity.ShiftAdmin{
				10: {ID: 10, TimeID: 2},
				30: {ID: 30, TimeID: 1},
			},
			wantIDs: []int{3, 1},
		},
		{
			name: "TimeIDゼロ値・負値は正値より前に並ぶ",
			logs: []entity.ActionLog{
				{ID: 1, ShiftID: 10},
				{ID: 2, ShiftID: 20},
				{ID: 3, ShiftID: 30},
			},
			shiftMap: map[int]entity.ShiftAdmin{
				10: {ID: 10, TimeID: 2},
				20: {ID: 20}, // TimeID はゼロ値 0
				30: {ID: 30, TimeID: -1},
			},
			wantIDs: []int{3, 2, 1},
		},
		{
			// 要判断（issue #TODO-B）: sort.Slice は非安定ソートのため、同一
			// TimeID 内の相対順序は仕様上未定義。ここでは要素が落ちない・
			// 重複しないことのみを集合一致で固定する。
			name: "同一TimeIDのタイは全件保持される（要判断: issue #TODO-B）",
			logs: []entity.ActionLog{
				{ID: 1, ShiftID: 10},
				{ID: 2, ShiftID: 20},
				{ID: 3, ShiftID: 30},
			},
			shiftMap: map[int]entity.ShiftAdmin{
				10: {ID: 10, TimeID: 1},
				20: {ID: 20, TimeID: 1},
				30: {ID: 30, TimeID: 1},
			},
			wantIDs:   []int{1, 2, 3},
			unordered: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := n.sortLogsByTime(c.logs, c.shiftMap)
			if got == nil {
				t.Fatalf("sortLogsByTime() = nil, want non-nil slice")
			}

			gotIDs := actionLogIDs(got)
			wantIDs := c.wantIDs
			if c.unordered {
				wantIDs = append([]int(nil), c.wantIDs...)
				sort.Ints(gotIDs)
				sort.Ints(wantIDs)
			}

			if !reflect.DeepEqual(gotIDs, wantIDs) {
				t.Errorf("sortLogsByTime() の ID 列 = %v, want %v", gotIDs, wantIDs)
			}
		})
	}

	t.Run("入力スライスを破壊しない", func(t *testing.T) {
		logs := []entity.ActionLog{
			{ID: 1, ShiftID: 10},
			{ID: 2, ShiftID: 20},
		}
		shiftMap := map[int]entity.ShiftAdmin{
			10: {ID: 10, TimeID: 2},
			20: {ID: 20, TimeID: 1},
		}

		got := n.sortLogsByTime(logs, shiftMap)

		if want := []int{2, 1}; !reflect.DeepEqual(actionLogIDs(got), want) {
			t.Errorf("sortLogsByTime() の ID 列 = %v, want %v", actionLogIDs(got), want)
		}
		if want := []int{1, 2}; !reflect.DeepEqual(actionLogIDs(logs), want) {
			t.Errorf("呼び出し後の入力の ID 列 = %v, want %v（入力が並べ替えられている）", actionLogIDs(logs), want)
		}
		if &got[0] == &logs[0] {
			t.Errorf("返り値が入力と同じ配列を指している（in-place ソートになっている）")
		}
	})
}

// テストの比較用に ActionLog スライスから ID 列を取り出す。
func actionLogIDs(logs []entity.ActionLog) []int {
	ids := make([]int, 0, len(logs))
	for _, log := range logs {
		ids = append(ids, log.ID)
	}
	return ids
}