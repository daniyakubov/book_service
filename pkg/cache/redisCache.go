package cache

import (
	"encoding/json"
	"time"

	"gopkg.in/redis.v5"
)

type RedisCache struct {
	host       string
	db         int
	expiration time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) BooksCache {
	return &RedisCache{
		host:       host,
		db:         db,
		expiration: exp,
	}
}
func (cache *RedisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *RedisCache) Set(key string, value Actions) {
	client := cache.getClient()
	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	client.Set(key, json, 0)
}

func (cache *RedisCache) Get(key string) *Actions {
	client := cache.getClient()
	val, err := client.Get(key).Result()
	if err != nil {
		return nil
	}
	var actions Actions
	err = json.Unmarshal([]byte(val), &actions)
	return &actions
}
