package cache

import (
	"context"
	"preacher61/go-assignment/model"
	"testing"

	"github.com/go-redis/redis/v8"
)

func getTestDataToSet() []*model.Activity {
	return []*model.Activity{
		{
			Activity: "test activity",
			Key:      "676767",
		},
		{
			Activity: "test activity-2",
			Key:      "676768",
		},
	}
}

func getTestRedisClient() *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	cli.FlushAll(context.Background())
	return cli
}

func TestRedisHandler(t *testing.T) {
	cli := getTestRedisClient()
	rcli := &RedisHandler{
		cli: cli,
	}

	ctx := context.Background()
	err := rcli.Set(ctx, getTestDataToSet())
	if err != nil {
		t.Fatal(err)
	}

	for v := range rcli.Get(ctx) {
		if len(v) != 1 {
			t.Fatalf("invalid value, expected lenght: 1, got: %d", len(v))
		}
	}
}
