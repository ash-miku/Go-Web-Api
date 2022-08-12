package router

import (
	"go-api/api"
	"go-api/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	//g.Use(common.NoCache)
	g.Use(common.Options)
	g.Use(common.Secure)
	g.Use(mw...)
	g.Use(common.RequestId)
	g.Use(common.SessionCheck)

	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// health check
	g.GET("/health", func(context *gin.Context) {
		common.LogInfo("check interfaces success")
	})

	g.GET("/v1/bili/video/info", api.GetVideoInfo)

	return g
}
