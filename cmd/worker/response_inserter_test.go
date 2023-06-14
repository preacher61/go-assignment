package main

import (
	"context"
	"preacher61/go-assignment/model"
	"testing"
)

func TestResponseInserterSuccess(t *testing.T) {
	ri := &responseInserter{
		iterateResponses: func(ctx context.Context) <-chan map[string][]*model.Activity {
			inChan := make(chan map[string][]*model.Activity, 1)
			defer close(inChan)

			inChan <- map[string][]*model.Activity{
				"responses_78787887": {
					{
						Activity: "test activity",
						Key:      "676767",
					},
					{
						Activity: "test activity-2",
						Key:      "676768",
					},
				},
			}
			return inChan
		},
		insert: func(ctx context.Context, data []*model.Activity) error {
			if len(data) != 2 {
				t.Fatal("len of data should be 2")
			}
			return nil
		},
		deleteKeys: func(ctx context.Context, keys []string) {
			if keys[0] != "responses_78787887" {
				t.Fatal("invalid key for deletion")
			}
		},
	}

	ri.processInsertion(context.Background())
}
