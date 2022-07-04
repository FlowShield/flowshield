package redis_hook

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

// HookConfig stores configuration needed to setup the hook
type HookConfig struct {
	Key      string
	Host     string
	Password string
	Port     int
	TTL      int
}

// RedisHook to sends logs to Redis server
type RedisHook struct {
	RedisPool      *redis.Pool
	RedisHost      string
	RedisKey       string
	LogstashFormat string
	AppName        string
	Hostname       string
	RedisPort      int
	TTL            int
}

// NewHook creates a hook to be added to an instance of logger
func NewHook(config HookConfig) (redisHook *RedisHook, err error) {
	pool := newRedisConnectionPool(config.Host, config.Password, config.Port, 0)

	// test if connection with REDIS can be established
	conn := pool.Get()
	defer conn.Close()

	// check connection
	_, err = conn.Do("PING")
	if err != nil {
		err = fmt.Errorf("unable to connect to REDIS: %s", err)
	}
	redisHook = &RedisHook{
		RedisHost:      config.Host,
		RedisPool:      pool,
		RedisKey:       config.Key,
		LogstashFormat: "origin",
		TTL:            config.TTL,
	}
	return
}

func newRedisConnectionPool(server, password string, port int, db int) *redis.Pool {
	hostPort := fmt.Sprintf("%s:%d", server, port)
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", hostPort, redis.DialDatabase(db),
				redis.DialPassword(password),
				redis.DialConnectTimeout(time.Second),
				redis.DialReadTimeout(time.Millisecond*100),
				redis.DialWriteTimeout(time.Millisecond*100))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
