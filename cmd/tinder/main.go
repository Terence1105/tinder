package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Terence1105/Tinder/cmd/httpserver"

	"github.com/Terence1105/Tinder/cmd/tinder/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		cancel()
	}()

	a := app.New(ctx)
	srv := httpserver.New(httpserver.ServeMux(a.Handler))
	srv.Execute(ctx)
}
