package tinder

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type TinderKV struct {
	conn *redis.Client
}

type TinderKVOption func(*TinderKV)

func WithRedisConn(conn *redis.Client) TinderKVOption {
	return func(kv *TinderKV) {
		kv.conn = conn
	}
}

func New(ctx context.Context, options ...TinderKVOption) *TinderKV {
	kv := &TinderKV{}
	for _, opt := range options {
		opt(kv)
	}

	return kv
}
