package backendforfrontend

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	redisClient *redis.Client
}

type RedisRepositoryI interface {
	FindByPostId(uuid.UUID) (PostAggregateI, error)
	UpdateCache(PostAggregateI)
}

func NewRedisRepository(redisClient *redis.Client) *RedisRepository {
	return &RedisRepository{
		redisClient: redisClient,
	}
}

func (r *RedisRepository) FindByPostId(postId uuid.UUID) (PostAggregateI, error) {
	postKey := fmt.Sprintf("post:%s", postId)
	postData, err := r.redisClient.Get(context.Background(), postKey).Result()
	if err != nil {
		return nil, err
	}

	var post PostAggregate
	if err := json.Unmarshal([]byte(postData), &post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *RedisRepository) UpdateCache(postAgg PostAggregateI) {
	postData, err := json.Marshal(postAgg)
	if err != nil {
		fmt.Println("Error during post aggregate serialization: " + err.Error())
		return
	}

	postKey := fmt.Sprintf("post:%s", postAgg.(*PostAggregate).Id)

	err = r.redisClient.Set(context.Background(), postKey, postData, 0).Err()
	if err != nil {
		fmt.Println("Something wrong with redis: " + err.Error())
	}
}
