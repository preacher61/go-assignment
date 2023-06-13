package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"preacher61/go-assignment/model"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

type activityAPI struct {
	rateLimiter *rate.Limiter
	url         string
}

func newActivityAPI(requests int, duration time.Duration) *activityAPI {
	return &activityAPI{
		rateLimiter: rate.NewLimiter(rate.Every(duration), requests),
		url:         "https://www.boredapi.com",
	}
}

func (a *activityAPI) getActivity(ctx context.Context) (*model.Activity, error) {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return nil, errors.Wrap(err, "api rate limit reached")
	}
	log.Info().Msg("activity_api: fetching.....")
	res, err := a.send(ctx)
	if err != nil {
		log.Error().Err(err).Msg("activity_api: send")
		return nil, errors.Wrap(err, "send")
	}
	log.Info().Interface("response", res).Msg("activity_api: response received")
	return res, nil
}

func (a *activityAPI) send(ctx context.Context) (*model.Activity, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/activity", a.url), http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, "new HTTP request")
	}

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "DO HTTP client request")
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.Errorf("got %d, want 2XX", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response")
	}

	var body *model.Activity
	err = json.Unmarshal(b, &body)
	if err != nil {
		return nil, errors.Wrap(err, "un-marshal response")
	}

	if body == nil {
		return nil, errors.New("empty response")
	}
	return body, nil
}
