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
// interval ごとに job を実行し、job が返したエラーはログ出力のみ（ループは止めない）。
func (s *Scheduler) Start(ctx context.Context) {
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
