package main

import (
	"github.com/NUTFes/seeft/seeft-api/api/lib/external/server"
	- "github.com/NUTFes/seeft/seeft-api/mysql"
)


func main() {
	db.Init()
	server.Init()
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
