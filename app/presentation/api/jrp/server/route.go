package server

import (
	"github.com/swaggo/echo-swagger"

	"github.com/yanosea/jrp/v2/app/presentation/api/jrp/server/jrp"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// Bind binds the routes to the server.
func Bind(e proxy.Echo) {
	e.Get("/swagger/*", echoSwagger.WrapHandler)
	apiGroup := e.Group("/api")
	jrp.BindGetJrpHandler(apiGroup)
}
