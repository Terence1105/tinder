package app

import (
	"runtime"
	"time"

	docs "github.com/Terence1105/Tinder/tool/swagger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Tinder API Doc
// @version 1.0
// @description This is a Tinder API.
// @host 127.0.0.1:8080
// @tag.name tinder
func NewGinEngine() *gin.Engine {
	engine := gin.New()
	// change to use zap logger
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "authorization", "X-API-Key"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	engine.GET("/health", Health)
	RegisterPrivateRouter(engine, true)
	docs.SwaggerInfo.BasePath = "/v1"
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return engine
}

func RegisterPrivateRouter(router *gin.Engine, enablePProf bool) {
	privateGroup := router.Group("/_")
	{
		if enablePProf {
			runtime.SetBlockProfileRate(1)
			runtime.SetMutexProfileFraction(1)
			pprof.RouteRegister(privateGroup, "/debug")
		}
	}
}

func (a *App) RegisterPublicRouter() {
	v1Group := a.Handler.Group("/v1")
	{
		v1Group.POST("/add-single-person-and-match", a.AddSinglePersonAndMatch)
		v1Group.POST("/remove-single-person", a.RemoveSinglePerson)
		v1Group.GET("/query-single-people", a.QuerySinglePeople)
	}
}

func Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
