package goengine

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
}

var ctx = context.Background()

func NewRedisStore(Addr string, Password string, DB int) SessionStore {
	client := redis.NewClient(&redis.Options{
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

func (rs *RedisStore) GetOnce(key string) (*map[string]interface{}, error) {
	store := make(map[string]interface{})
	jsonSession, err := rs.client.Eval(
		ctx,
		"local ret = redis.call('GET', KEYS[1]); redis.call('DEL', KEYS[1]); return ret",
		[]string{key},
	).Text()

	if nil == err {
		json.Unmarshal([]byte(jsonSession), &store)
	}

	if redis.Nil == err {
		err = nil
	}

	return &store, err
}

func (rs *RedisStore) Get(key string) (*map[string]string, error) {
	store := make(map[string]string)
	jsonSession, err := rs.client.Get(ctx, key).Result()

	if nil == err {
		json.Unmarshal([]byte(jsonSession), &store)
	}

	if redis.Nil == err {
		err = nil
	}

	return &store, err
}

/**
 * @param json session data in json format, if nil or empty, delete the key
 * @param maxAge in seconds, if maxAge ==0, no expiration; if maxAge <0, delete the key
 */
func (rs *RedisStore) Save(key string, json []byte, maxAge int) error {
	if maxAge < 0 || nil == json || 0 == len(json) {
		return rs.client.Del(ctx, key).Err()
	}
	return rs.client.Set(ctx, key, string(json), time.Duration(maxAge)*time.Second).Err()
}
