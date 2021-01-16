package redisdb

import (
	"fmt"
	"github.com/phpunch/route-roam-api/log"
	"time"

	"github.com/gomodule/redigo/redis"
)

type DB interface {
	Set(key string, value interface{}, expDuration time.Duration) error
	Get(key string) (int64, error)
	Del(key string) (int64, error)
}

type redisDB struct {
	pool *redis.Pool
}

func NewPool(config *Config) *redis.Pool {
	return &redis.Pool{
		MaxActive: 4000,
		MaxIdle:   100,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Address)
			if err != nil {
				log.Log.Fatalf(err.Error())
			}
			return c, err
		},
	}
}

func New(config *Config) DB {
	pool := NewPool(config)
	return &redisDB{pool: pool}
}

func (r *redisDB) Set(key string, value interface{}, expDuration time.Duration) error {
	// jsonData, err := json.Marshal(value)
	// if err != nil {
	// 	return fmt.Errorf("could not marshal json : %v", err)
	// }
	secs := int(expDuration.Seconds())
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", key, secs, value)
	if err != nil {
		return fmt.Errorf("could not set data : %v", err)
	}

	return nil
}

func (r *redisDB) Get(key string) (int64, error) {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.Int64(conn.Do("GET", key))
	if err != nil {
		return 0, fmt.Errorf("could not get data : %v", err)
	}
	return val, nil
}

func (r *redisDB) Del(key string) (int64, error) {
	conn := r.pool.Get()
	defer conn.Close()
	val, err := redis.Int64(conn.Do("DEL", key))
	if err != nil {
		return 0, fmt.Errorf("could not del data : %v", err)
	}
	return val, nil
}
