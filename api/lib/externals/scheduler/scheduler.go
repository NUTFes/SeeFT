package scheduler

import (
	"context"
	"log"
	"time"
)

// Job は scheduler が定期実行する処理。usecase を import せず関数型で受け取り、
// scheduler と業務ロジックを疎結合に保つ。
type Job func(ctx context.Context) error

// Scheduler は「名前・間隔・実行する処理」を保持する。
type Scheduler struct {
	name     string
	interval time.Duration
	job      Job
}

// New はコンストラクタでSchedulerを生成する。
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
