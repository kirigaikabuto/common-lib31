package common

import (
	"encoding/json"
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

func (r *RedisStore) SetValue(key string, value interface{}, t time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = r.client.Set(key, jsonData, t).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisStore) GetValue(key string, src interface{}) error {
	val, err := r.client.Get(key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), src)
	if err != nil {
		return err
	}
	return nil
}
