package main

import (
	"github.com/NUTFes/SeeFT/api/lib/internals/di"
	_ "github.com/go-sql-driver/mysql"
)


func main() {
	client := di.InitializeServer()
	defer client.CloseDB()
	// db.Init()
	// server.Init()
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
