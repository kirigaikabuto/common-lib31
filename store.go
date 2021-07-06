package common

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisConnect(cf RedisConfig) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cf.Host + ":" + cf.Port,
	})
	err := client.Ping().Err()
	if err != nil {
		return nil, err
	}
	return &RedisStore{client: client}, nil
}

func (r *RedisStore) Save(key string, value interface{}, t time.Duration) error {
	err := r.client.Set(key, value, t).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisStore) GetValue(key string) (string, error) {
	val, err := r.client.Get(key).Result()
	if err != nil {
		return "", nil
	}
	return val, nil
}
