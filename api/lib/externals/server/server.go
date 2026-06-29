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

func RunServer(ctx context.Context, router router.Router) error {
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

	errCh := make(chan error, 1)

	// サーバー起動
	go func() {

		if err := e.Start(":1234"); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():

	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return e.Shutdown(shutdownCtx)
}
