package redisConn

import (
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

var Connection *RedisConnection

type RedisConnection struct {
	client *redis.Client
}

func NewRedisConnection(addr, password string, db int) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	Connection = &RedisConnection{client: client}
}

func (r *RedisConnection) SetToken(token string) error {
	err := r.client.Set(token, true, time.Hour).Err()
	return err
}

func (r *RedisConnection) GetToken(token string) (bool, error) {
	val, err := r.client.Get(token).Result()
	if err != nil {
		return false, err
	}

	res, err := strconv.ParseBool(val)
	return res, err
}
