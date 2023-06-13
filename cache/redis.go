package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"preacher61/go-assignment/model"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var cacheExpiry = 24 * time.Hour

// RedisCache encapsulates redis client.
type RedisClient struct {
	cli *redis.Client
}

// NewRedisClient returns a new redis client.
func NewRedisClient() *RedisClient {
	cli := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisClient{
		cli: cli,
	}
}

func (r *RedisClient) Set(ctx context.Context, data []*model.Activity) error {
	log.Info().Interface("data", data).Msg("redis: inserting responses")

	b, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}
	err = r.cli.Set(ctx, r.getKey(), b, cacheExpiry).Err()
	if err != nil {
		log.Error().Err(err).Msg("redis: error occured while inserting")
		return errors.Wrap(err, "set")
	}
	log.Info().Msg("redis: responses inserted")
	return nil
}

func (r *RedisClient) getKey() string {
	return fmt.Sprintf("response_%d", time.Now().Unix())
}
