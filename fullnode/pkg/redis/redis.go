package redis

import (
	"context"

	"github.com/flowshield/flowshield/fullnode/pkg/confer"
	"github.com/go-redis/redis/v8"
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
