// scheduler.Start の挙動を目で見るための使い捨てデモ。確認後に削除する。
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/NUTFes/SeeFT/api/lib/externals/scheduler"
)

func main() {
	fmt.Println("[1] Start を呼ぶ前")

	// 本番では NotificationUseCase.ProcessUnsentNotifications が入る。
	// デモでは Slack も DB も使わず、print だけするダミー job。
	job := func(ctx context.Context) error {
		fmt.Println("    → job 実行！（本番ではここで Slack DM を送る）")
		return nil
	}

	// 間隔は 2 秒（本番は 5*time.Minute）。これが s.interval になる。
	scheduler.New("demo", 2*time.Second, job).Start(context.Background())

	fmt.Println("[2] Start を呼んだ直後（もうここに来た＝ブロックしてない）")

	// 本番では server.RunServer がここでブロックしてプロセスを生かし続ける。
	// デモでは 7 秒だけ待って、裏で job が何回鳴るか観察する。
	time.Sleep(7 * time.Second)
	fmt.Println("[3] デモ終了（main が終わると裏方の goroutine も道連れに消える）")
}
