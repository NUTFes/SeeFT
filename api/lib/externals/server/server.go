package server

import (
	"context"
	"net/http"
	"os"
	"time"

	_ "github.com/NUTFes/SeeFT/api/docs"
	"github.com/NUTFes/SeeFT/api/lib/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RunServer(ctx context.Context, router router.Router) {
	// echoのインスタンス
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())

	// log
	logger := middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: logFormat(),
			Output: os.Stdout,
		},
	)
	e.Use(logger)

	// CORS対策
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "127.0.0.1:3000", "http://localhost:3001", "127.0.0.1:3001", "http://localhost:45029", "127.0.0.1:45029", "https://seeft.nutfes.net", "https://seeft-admin.nutfes.net"}, // ドメイン
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// ルーティング
	router.ProvideRouter(e)

	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// サーバー起動
	go func() {

		if err := e.Start(":1234"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}
}
