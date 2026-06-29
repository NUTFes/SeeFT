package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/NUTFes/SeeFT/api/lib/di"
	_ "github.com/lib/pq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	client := di.InitializeServer(ctx)
	defer client.CloseDB()
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
