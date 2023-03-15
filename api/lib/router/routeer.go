package router

import (
	"github.com/NUTFes/SeeFT/api/lib/externals/controller"
	"github.com/labstack/echo/v4"
)

type router struct {
	healthcheckController     controller.HealthcheckController
	bureauController          controller.BureauController
}

type Router interface {
	ProvideRouter(*echo.Echo)
}

func NewRouter(
	healthController controller.HealthcheckController,
	bureauController controller.BureauController,

) Router {
	return router{
		healthController,
		bureauController,
	}
}

func (r router) ProvideRouter(e *echo.Echo) {
	// Healthcheck
	e.GET("/", r.healthcheckController.IndexHealthcheck)

	//bureauのRoute
	e.GET("/bureaus", r.bureauController.IndexBureau)
	e.GET("/bureaus/:id", r.bureauController.ShowBureau)
}