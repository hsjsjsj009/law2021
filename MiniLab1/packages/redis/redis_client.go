package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type Client struct {
	Client *redis.Client
}

func (redisClient Client) StoreToRedisWithExpired(key string, val interface{}, duration string) error {
	dur, err := time.ParseDuration(duration)
	if err != nil {
		return err
	}

	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = redisClient.Client.Set(key, string(b), dur).Err()

	return err
}


func (redisClient Client) StoreToRedis(key string, val interface{}) error {
	b, err := json.Marshal(val)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = redisClient.Client.Set(key, string(b), 0).Err()

	return err
}

func (redisClient Client) GetFromRedis(key string, cb interface{}) error {
	res, err := redisClient.Client.Get(key).Result()
	if err != nil {
		return err
	}

	if res == "" {
		return errors.New("[Redis] Value of " + key + " is empty.")
	}

	err = json.Unmarshal([]byte(res), cb)
	if err != nil {
		return err
	}

	return err
}

func (redisClient Client) RemoveFromRedis(key string) error {
	return redisClient.Client.Del(key).Err()
}