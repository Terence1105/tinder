package app

import (
	"context"

	"github.com/Terence1105/Tinder/pkg/storage"
	"github.com/gin-gonic/gin"
)

type App struct {
	Handler *gin.Engine
	storage storage.TinderStorage
}

func New(ctx context.Context, storage storage.TinderStorage) *App {
	app := &App{
		Handler: NewGinEngine(),
		storage: storage,
	}

	app.RegisterPublicRouter()

	return app
}
