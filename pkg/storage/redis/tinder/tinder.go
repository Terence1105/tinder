package tinder

import (
	"context"
	"fmt"

	"github.com/Terence1105/Tinder/pkg/storage/redis/tinder/dto"
	"github.com/Terence1105/Tinder/pkg/types"
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

func (t *TinderKV) AddPerson(ctx context.Context, person *dto.Person) error {
	var heightKey string
	if person.Gender == types.BOY {
		heightKey = KEY_BOYS_HEIGHT
	} else {
		heightKey = KEY_GIRLS_HEIGHT
	}

	zadd := t.conn.ZAdd(ctx, heightKey, &redis.Z{
		Score:  person.Height,
		Member: person.Name,
	})
	if err := zadd.Err(); err != nil {
		return err
	}

	dateCountsKey := fmt.Sprintf(KEY_USER_DATE_COUNTS, person.Name)
	set := t.conn.Set(ctx, dateCountsKey, person.DateCounts, 0)
	if err := set.Err(); err != nil {
		return err
	}

	return nil
}

func (t *TinderKV) RemovePerson(ctx context.Context, name string, gender int) error {
	var heightKey string
	if gender == types.BOY {
		heightKey = KEY_BOYS_HEIGHT
	} else {
		heightKey = KEY_GIRLS_HEIGHT
	}

	zrem := t.conn.ZRem(ctx, heightKey, name)
	if err := zrem.Err(); err != nil {
		return err
	}

	dateCountsKey := fmt.Sprintf(KEY_USER_DATE_COUNTS, name)
	del := t.conn.Del(ctx, dateCountsKey)
	if err := del.Err(); err != nil {
		return err
	}

	return nil
}

func (t *TinderKV) GetPeople(ctx context.Context, min, max float64, count, gender int) ([]dto.Person, error) {
	var key string

	if gender == types.BOY {
		key = KEY_BOYS_HEIGHT
	} else {
		key = KEY_GIRLS_HEIGHT
	}

	users, err := t.conn.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min:    fmt.Sprintf("(%f", min),
		Max:    fmt.Sprintf("(%f", max),
		Offset: 0,
		Count:  int64(count),
	}).Result()

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}
	resp := []dto.Person{}

	for _, v := range users {
		boy := dto.Person{
			Name:   v.Member.(string),
			Height: v.Score,
			Gender: gender,
		}
		resp = append(resp, boy)
	}

	return resp, nil
}

func (t *TinderKV) DecrementDateCount(ctx context.Context, name string) (int, error) {
	dateCountsKey := fmt.Sprintf(KEY_USER_DATE_COUNTS, name)

	decr := t.conn.Decr(ctx, dateCountsKey)
	if err := decr.Err(); err != nil {
		return -1, err
	}

	return int(decr.Val()), nil
}

func (t *TinderKV) GetDateCount(ctx context.Context, name string) (string, error) {
	dateCountsKey := fmt.Sprintf(KEY_USER_DATE_COUNTS, name)

	get := t.conn.Get(ctx, dateCountsKey)
	if err := get.Err(); err != nil {
		return "", err
	}

	return get.Val(), nil
}
