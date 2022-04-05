package cache

import (
	"time"

	errors "github.com/fiverr/go_errors"
	"gopkg.in/redis.v5"
)

var _ Cache = &RedisCache{}

type RedisCache struct {
	host       string
	db         int
	expiration time.Duration
	maxSize    int64
	client     *redis.Client
}

func NewRedisCache(host string, db int, exp time.Duration, maxSize int64, client *redis.Client) *RedisCache {
	return &RedisCache{
		host:       host,
		db:         db,
		expiration: exp,
		maxSize:    maxSize,
		client:     client,
	}
}
func (cache *RedisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *RedisCache) Push(key string, value string) error {
	length, err := cache.client.LLen(key).Result()
	if err != nil {
		panic(err)
	}

	cache.client.RPush(key, value)

	if length >= cache.maxSize {
		_, err := cache.client.LPop(key).Result()
		if err != nil {
			return errors.Wrap(err, err.Error())
		}
	}
	return nil
}

func (cache *RedisCache) Get(key string) ([]string, error) {
	length, err := cache.client.LLen(key).Result()
	if err != nil {
		panic(err)
	}
	actions, err := cache.client.LRange(key, 0, length).Result()
	if err != nil {
		return nil, errors.Wrap(err, err.Error())
	}
	return actions, nil
}
