package server

import (
	"github.com/yanosea/jrp/v2/app/presentation/api/jrp/server/jrp"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// Bind binds the routes to the server.
func Bind(e proxy.Echo) {
	apiGroup := e.Group("/api")
	jrp.BindGetJrpHandler(apiGroup)
}
