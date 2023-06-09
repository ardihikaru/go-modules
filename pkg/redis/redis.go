package redis

import (
	"time"

	"encoding/json"
	"github.com/go-redis/redis"
)

type Redis struct {
	Client *redis.Client
}

// GetRedis build redis object
func GetRedis(addr, passwd string, db int) (*Redis, error) {
	redisClient, err := Connect(addr, passwd, db)
	if err != nil {
		return nil, err
	}

	// builds redis object
	redisObj := &Redis{
		Client: redisClient,
	}

	return redisObj, nil
}

// Connect connects redis database
func Connect(addr, passwd string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Set sets value to redis database
func (r *Redis) Set(key string, value interface{}, exp time.Duration) error {
	var err error

	p, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.Client.Set(key, p, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

// Get gets value from redis database
func (r *Redis) Get(key string, destType interface{}) error {
	val, err := r.Client.Get(key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), &destType)
}
