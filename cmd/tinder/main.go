package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Terence1105/Tinder/cmd/httpserver"
	"github.com/Terence1105/Tinder/pkg/storage/redis/tinder"
	"github.com/go-redis/redis/v8"

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

	// TODO: use config
	redis, err := ConnKevVal(ctx, "redis:6379", "", 0, 10)
	if err != nil {
		panic(err)
	}
	a := app.New(ctx, redis)

	srv := httpserver.New(httpserver.ServeMux(a.Handler))
	srv.Execute(ctx)
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
