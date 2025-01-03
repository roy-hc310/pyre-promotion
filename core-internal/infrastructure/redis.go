package infrastructure

import (
	"pyre-promotion/core-internal/utils"
	"time"

	"github.com/go-redis/redis"
	// "github.com/rs/zerolog/log"
)

type RedisInfra struct {
	Client *redis.Client
}

func NewRedisInfra() *RedisInfra {

	client := redis.NewClient(&redis.Options{
		Addr:        utils.GlobalEnv.Redis.Host,
		Password:    utils.GlobalEnv.Redis.Pass,
		DB:          0,
		DialTimeout: time.Second * time.Duration(utils.GlobalEnv.Redis.TimeOut),
		MaxRetries:  0,
	})

	return &RedisInfra{
		Client: client,
	}
}
