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

// RedisHandler encapsulates redis client.
type RedisHandler struct {
	cli *redis.Client
}

// NewRedisHandler returns a new redis client.
func NewRedisHandler() *RedisHandler {
	cli := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisHandler{
		cli: cli,
	}
}

func (r *RedisHandler) Set(ctx context.Context, data []*model.Activity) error {
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

func (r *RedisHandler) getKey() string {
	return fmt.Sprintf("response_%d", time.Now().Unix())
}

// Get traverses over keys returning values corresponding to it.
func (r *RedisHandler) Get(ctx context.Context) <-chan map[string][]*model.Activity {
	inStream := make(chan map[string][]*model.Activity, 1)

	iter := r.cli.Scan(ctx, 0, "response_*", 0).Iterator()
	go func() {
		defer close(inStream)

		for iter.Next(ctx) {
			select {
			case <-ctx.Done():
				return
			default:
				key := iter.Val()

				val, err := r.cli.Get(ctx, key).Result()
				if err != nil {
					log.Error().Err(err).Msg("get redis key failed")
					continue
				}

				var v []*model.Activity
				err = json.Unmarshal([]byte(val), &v)
				if err != nil {
					log.Error().Err(err).Msg("parsing redis response failed")
					continue
				}
				inStream <- map[string][]*model.Activity{
					key: v,
				}
			}
		}

	}()

	return inStream
}

// DeleteMulti deletes multiple keys passed in as an argument.
func (r *RedisHandler) DeleteMulti(ctx context.Context, keys []string) {
	log.Info().Msg("deleting keys")

	pipe := r.cli.Pipeline()
	pipe.Del(ctx, keys...)
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error occured while deleting redis keys")
	}
}
