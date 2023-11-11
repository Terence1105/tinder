package app

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/Terence1105/Tinder/pkg/storage/redis/tinder"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

type App struct {
	Handler *gin.Engine
	// TODO: use interface
	redisCli *tinder.TinderKV
}

func New(ctx context.Context) *App {
	redis, err := ConnKevVal(ctx, "redis:6379", "", 0, 10)
	if err != nil {
		panic(err)
	}

	app := &App{
		Handler:  NewGinEngine(),
		redisCli: redis,
	}

	app.RegisterPublicRouter()

	return app
}

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
		// TODO: 
		// prometheus.Register(_RouterMetrics)
		privateGroup.GET("/metrics", gin.WrapH(promhttp.Handler()))
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

func ConnKevVal(ctx context.Context, addr string, password string, db, poolSize int) (*tinder.TinderKV, error) {
	_, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	opt := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	}

	conn := redis.NewClient(opt)

	// check
	ct, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err := conn.Ping(ct).Result()
	if err != nil {
		return nil, fmt.Errorf("conn redis fail : %w", err)
	}

	var opts []tinder.TinderKVOption
	opts = append(opts, tinder.WithRedisConn(conn))
	kv := tinder.New(ctx, opts...)

	return kv, nil
}
