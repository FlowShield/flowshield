package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/flowshield/flowshield/fullnode/pkg/confer"
)

var Client *redis.Client

func Init(cfg *confer.Redis) (err error) {
	Client = redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
	if err = Client.Ping(context.Background()).Err(); err != nil {
		return
	}
	return
}
