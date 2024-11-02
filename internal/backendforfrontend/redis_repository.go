package backendforfrontend

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	redisClient *redis.Client
}

type RedisRepositoryI interface {
	UpdateCache(post PostAggregateI)
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	return &RedisRepository{
		redisClient: redisClient,
	}
}

func (r *RedisRepository) UpdateCache(postAgg PostAggregateI) {
	postData, err := json.Marshal(postAgg)
	if err != nil {
		fmt.Printf("Error during post aggregate serialization: " + err.Error())
	}

	postKey := fmt.Sprintf("post:%s", postAgg.(*PostAggregate).Id)

	err = r.redisClient.Set(context.Background(), postKey, postData, 0).Err()
	if err != nil {
		fmt.Printf("Something wrong with redis: " + err.Error())
	}
}
