package redisdb

import (
	"encoding/json"
	"fmt"
	"github.com/phpunch/route-roam-api/log"
	"time"

	"github.com/gomodule/redigo/redis"
)

type DB interface {
	Set(key string, value interface{}, expDuration time.Duration) error
	Get(key string) ([]map[string]interface{}, error)
}

type redisDB struct {
	Conn redis.Conn
}

func NewPool(config *Config) *redis.Pool {
	return &redis.Pool{
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
	conn := pool.Get()
	return &redisDB{Conn: conn}
}

func (r *redisDB) Set(key string, value interface{}, expDuration time.Duration) error {
	// jsonData, err := json.Marshal(value)
	// if err != nil {
	// 	return fmt.Errorf("could not marshal json : %v", err)
	// }
	secs := int(expDuration.Seconds())
	_, err := r.Conn.Do("SETEX", key, secs, value)
	if err != nil {
		return fmt.Errorf("could not set data : %v", err)
	}

	return nil
}

func (r *redisDB) Get(key string) ([]map[string]interface{}, error) {
	val, err := redis.String(r.Conn.Do("GET", key))
	if err != nil {
		return nil, fmt.Errorf("could not get data : %v", err)
	}
	var result []map[string]interface{}
	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json : %v", err)
	}

	return result, nil
}
