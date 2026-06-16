package scheduler

import (
	"context"
	"log"
	"time"
)

// Job は定期実行する処理の型。
//
//	Go:     type Job func(ctx context.Context) error
//	Python: Callable[[Context], Awaitable[None]]  （typing の型エイリアス）
//	TS:     type Job = (ctx: Context) => Promise<void>
//
// 関数を「型」として名付けることで、scheduler は usecase を import せずに済む
// （依存方向が di → scheduler / di → usecase の二股になる）。
type Job func(ctx context.Context) error

// Scheduler は「名前・間隔・実行する処理」を保持するだけの箱。
// Python なら @dataclass、TS なら class { constructor(...) } に相当。
type Scheduler struct {
	name     string
	interval time.Duration
	job      Job
}

// New はコンストラクタ。Go には __init__ / constructor が無いので、
// 慣習として New〜 関数でポインタを返す。
func New(name string, interval time.Duration, job Job) *Scheduler {
	return &Scheduler{name: name, interval: interval, job: job}
}

// Start は ticker ループを goroutine で起動し、即座に return する。
//
// ★ここの中身は上林さんが書く。
//   - go func() { ... }() で別の実行単位（goroutine）を起動して呼び出し元をブロックしない
//   - time.NewTicker(s.interval) で s.interval ごとに発火するチャネルを作る
//   - for ループ内で <-ticker.C を待ち、発火のたびに s.job(ctx) を呼ぶ
//   - job がエラーを返したら log.Printf でログ出力のみ（パニックさせない）
//
// 実装したら "log" を import に追加すること（Go は未使用 import をコンパイルエラーにする）。
func (s *Scheduler) Start(ctx context.Context) {
	// TODO(上林): goroutine + time.NewTicker ループを実装する
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for range ticker.C {
			if err := s.job(ctx); err != nil {
				log.Printf("[scheduler:%s] job error: %v", s.name, err)
			}
		}
	}()
}
