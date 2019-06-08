package goengine

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

type RedisStore struct {
	client *redis.Client
}

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
	jsonSession, err := this.client.Get(key).Result()

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
	return this.client.Set(key, string(buf), time.Duration(maxAge)).Err()
}
