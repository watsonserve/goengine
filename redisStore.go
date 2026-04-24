package goengine

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
}

var ctx = context.Background()

func NewRedisStore(Addr string, Password string, DB int) *RedisStore {
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

func (rs *RedisStore) GetOnce(key string, store interface{}) error {
	jsonSession, err := rs.client.Eval(
		ctx,
		"local ret = redis.call('GET', KEYS[1]); redis.call('DEL', KEYS[1]); return ret",
		[]string{key},
	).Text()

	if nil == err {
		json.Unmarshal([]byte(jsonSession), store)
	}

	if redis.Nil == err {
		err = nil
	}

	return err
}

func (rs *RedisStore) Get(key string, store interface{}) error {
	jsonSession, err := rs.client.Get(ctx, key).Result()

	if nil == err {
		json.Unmarshal([]byte(jsonSession), store)
	}

	if redis.Nil == err {
		err = nil
	}

	return err
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

func (rs *RedisStore) Expire(key string, maxAge int) error {
	if maxAge < 1 {
		return errors.New("maxAge must be greater than 0")
	}
	return rs.client.Expire(ctx, key, time.Duration(maxAge)*time.Second).Err()
}
