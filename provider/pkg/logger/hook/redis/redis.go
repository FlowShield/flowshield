package redis

import (
	"fmt"

	"github.com/cloudslit/cloudslit/provider/pkg/util/json"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Addr     string
	Password string
	Key      string
}

func New(c *Config) *Hook {
	redisdb := redis.NewClient(&redis.Options{
		Addr:     c.Addr,     // use default Addr
		Password: c.Password, // no password set
		DB:       0,          // use default DB
	})

	_, err := redisdb.Ping().Result()
	if err != nil {
		fmt.Println("error creating message for REDIS:", err)
		panic(err)
	}
	return &Hook{
		cli: redisdb,
		key: c.Key,
	}
}

type Hook struct {
	cli *redis.Client
	key string
}

func (h *Hook) Exec(entry *logrus.Entry) error {
	fields := make(map[string]interface{})
	for k, v := range entry.Data {
		fields[k] = v
	}
	fields["level"] = entry.Level.String()
	fields["message"] = entry.Message
	b, _ := json.Marshal(fields)
	return h.cli.RPush(h.key, string(b)).Err()
}

func (h *Hook) Close() error {
	return h.cli.Close()
}
