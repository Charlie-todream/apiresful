package router

import (
	"net/http"

	"apiresful/handler/sd"
	"apiresful/handler/user"
	"apiresful/router/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	_ "apiresful/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	pprof.Register(g)
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	u := g.Group("/v1/user")
	u.Use(middleware.AuthMiddleware())
	{
		u.POST("", user.Create)

		u.DELETE("/:id", user.Delete)
		u.PUT("/:id", user.Update)
		u.GET("", user.List)
		u.GET("/:username", user.Get)

	}
	g.POST("/login", user.Login)

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return g
}