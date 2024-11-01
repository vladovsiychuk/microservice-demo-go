package backendforfrontend

import (
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	redisClient *redis.Client
}

type RedisRepositoryI interface {
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	return &RedisRepository{}
}
