package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NUTFes/SeeFT/api/lib/di"
	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	client, err := di.InitializeServer(ctx)
	if client != nil {
		defer client.CloseDB()
	}
	return err
}

// func main() async {
//   final Environment env = Environment();
//   if (env.applicationEnv == "production") {
//     Log.setupProd();
//   } else {
//     Log.setupDev();
//   }

//   await HotReloader.create(
//       onAfterReload: (ctx) => logger.info("Hot-reload result: ${ctx.result}\n ${ctx.reloadReports}"));

//   final server = await initializeServer(env);

//   await server.run();
// }
