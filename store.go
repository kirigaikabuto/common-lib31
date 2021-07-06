package common

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisConnectStore struct {
	client *redis.Client
}

func NewRedisConnectStore(cfg RedisConfig) (*RedisConnectStore, error) {
	//connection to redis
	client := redis.NewClient(&redis.Options{Addr: cfg.Host + ":" + cfg.Port})
	err := client.Ping().Err()
	if err != nil {
		return nil, err
	}
	return &RedisConnectStore{client: client}, nil
}

func (r *RedisConnectStore) Save(key string, value interface{}, t time.Duration) error {
	//convert to json if all values is only json
	err := r.client.Set(key, value, t).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisConnectStore) Get(key string) (string, error) {
	value, err := r.client.Get(key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
