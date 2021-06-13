package goengine

import (
	"context"
	"encoding/json"
	redis "github.com/go-redis/redis/v8"
	"time"
)

type RedisStore struct {
	client *redis.Client
}

var ctx = context.Background()

func NewRedisStore(Addr string, Password string, DB int) *RedisStore {
	client := redis.NewClient(&redis.Options {
		Addr:     Addr,
		Password: Password, // no password set
		DB:       DB,       // use default DB
	})
	if nil == client {
		return nil
	}
	return &RedisStore{
		client: client,
	}
}

func (this *RedisStore) Get(key string) (*map[string]interface{}, error) {
	store := make(map[string]interface{})
	jsonSession, err := this.client.Get(ctx, key).Result()

  if nil == err {
		json.Unmarshal([]byte(jsonSession), &store)
	}

	if redis.Nil == err {
		err = nil
	}

	return &store, err
}

func (this *RedisStore) Save(key string, store *map[string]interface{}, maxAge int) error {
	buf, err := json.Marshal(store)
	if nil != err {
		return err
	}
	return this.client.Set(ctx, key, string(buf), time.Duration(maxAge) * time.Second).Err()
}
