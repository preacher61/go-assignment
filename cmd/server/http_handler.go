package main

import (
	"context"
	"log"
	"net/http"
	"preacher61/go-assignment/httpjson"
	"preacher61/go-assignment/model"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type httpGetEventsHandler struct {
	fetchActivity func(ctx context.Context) (*model.Activity, error)
}

func newHTTPGetEventsHandler() *httpGetEventsHandler {
	a := newActivityAPI(3, time.Second)
	return &httpGetEventsHandler{
		fetchActivity: a.getActivity,
	}
}

func (h *httpGetEventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("request received")
	ctx := r.Context()
	res, err := h.handle(ctx)
	if err != nil {
		err = errors.Wrap(err, "handle")
		httpjson.WriteResponse(w, http.StatusInternalServerError, &errorResponse{
			Error: err.Error(),
		})
		return
	}
	httpjson.WriteResponse(w, http.StatusOK, res)
}

func (h *httpGetEventsHandler) handle(ctx context.Context) ([]*model.Activity, error) {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	respChan := make(chan *model.Activity)
	errChan := make(chan error)
	done := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)

		go func(ctx context.Context) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				res, err := h.fetchActivity(ctx)
				if err != nil {
					err = errors.Wrap(err, "fetch activity")
					errChan <- err
					return
				}
				respChan <- res
			}
		}(cancelCtx)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	responses := []*model.Activity{}
	for {
		select {
		case err := <-errChan:
			cancel()
			return nil, err
		case res := <-respChan:
			responses = append(responses, res)
		case <-done:
			return responses, nil
		case <-time.After(2 * time.Second):
			cancel()
			return nil, errors.New("Activity-API not available")
		}
	}
}
