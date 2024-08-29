package infrastructure

import (
	"pyre-promotion/core-internal/utils"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	// "github.com/rs/zerolog/log"
)

type RedisInfra struct {
	Client *redis.Client
}

func NewRedisInfra() *RedisInfra {

	timeOut, _ := strconv.Atoi(utils.GlobalEnv.Redis.TimeOut)
	client := redis.NewClient(&redis.Options{
		Addr:        utils.GlobalEnv.Redis.HostPort,
		Password:    utils.GlobalEnv.Redis.Pass,
		DB:          0,
		DialTimeout: time.Second * time.Duration(timeOut),
		MaxRetries:  0,
	})

	return &RedisInfra{
		Client: client,
	}
}
