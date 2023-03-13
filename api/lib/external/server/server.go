package server

import (
	"github.com/NUTFes/seeft-api/lib/interface/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init() {
	r := router()
	r.Run(":3000")
}

func router() *gin.Engine {
	r := gin.Default()

	r.Use(CORS())

	health := r.Group("/healthz")
	{
		healthController := controller.HealthController{}
		health.GET("/", healthController.GetHealth)
	}

	user := r.Group("/user")
	{
		userController := controller.UserController{}
		user.GET("/:id", userController.FindByID)
		user.GET("/", userController.Index)
	}

	auth := r.Group("/auth")
	{
		authController := controller.AuthController{}
		auth.POST("/", authController.Index)
		auth.POST("/attendance", authController.PostCheck)
		auth.GET("/:mail", authController.MailAuth)
	}

	shift := r.Group("/shift")
	{
		shiftController := controller.ShiftController{}
		shift.GET("/:id", shiftController.Search)
		shift.GET("/:id/:day", shiftController.Search)
		shift.GET("/:id/:day/:weather", shiftController.Search)
	}

	work := r.Group("/work")
	{
		workController := controller.WorkController{}
		work.GET("/:workID/:userID/:day/:weather/:time", workController.WorkDetail)
		work.GET("/list", workController.WorkList)
	}

	return r
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// import 'package:shelf/shelf.dart';
// import 'package:shelf/shelf_io.dart' as shelf_io;

// import 'service.dart';
// import '../../config/config.dart';

// class Server {
//   final Service service;
//   final Environment env;

//   Server(this.service, this.env);

//   run() async {
//     final handler = _handler();
//     final server = await shelf_io.serve(
//       handler,
//       env.server.address,
//       int.parse(env.server.port),
//     );
//     print('Server runnning on localhost:${server.port}');
//   }

//   _handler() {
//     final pipeline = Pipeline().addMiddleware(logRequests()).addHandler(service.handler);
//     return pipeline;
//   }
// }
