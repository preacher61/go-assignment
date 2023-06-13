package main

import (
	"context"
	"net/http"
	"preacher61/go-assignment/cache"
	"preacher61/go-assignment/httpjson"
	"preacher61/go-assignment/model"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var errActivityApiNotAvailabe = errors.New("Activity-API not available")

type httpGetEventsHandler struct {
	fetchActivity   func(ctx context.Context) (*model.Activity, error)
	persistResponse func(ctx context.Context, response []*model.Activity) error
}

func newHTTPGetEventsHandler() *httpGetEventsHandler {
	a := newActivityAPI(3, time.Second)
	rCli := cache.NewRedisClient()
	return &httpGetEventsHandler{
		fetchActivity:   a.getActivity,
		persistResponse: rCli.Set,
	}
}

func (h *httpGetEventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := h.handle(ctx)
	if err != nil {
		err = errors.Wrap(err, "handle")
		httpjson.WriteResponse(w, http.StatusInternalServerError, &errorResponse{
			Error: err.Error(),
		})
		log.Error().Err(err).Msg("serve HTTP")
		return
	}
	log.Info().Msg("success")
	httpjson.WriteResponse(w, http.StatusOK, res)
}

type responseMeta struct {
	response *model.Activity
	err      error
}

func (h *httpGetEventsHandler) handle(ctx context.Context) ([]*model.Activity, error) {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	responseChan := make(chan []*model.Activity)
	responseChanToPersist := make(chan []*model.Activity) // a chan that would help us persist the resposnes.
	errChan := make(chan error)
	persistingDone := make(chan struct{}) // notifies if the persisting process has been completed, if yes we can return the resposne.

	// call the APIs and prepare the responses.
	go h.prepareResponse(cancelCtx, responseChan, errChan)

	// persist the above responses into some cache/messaging queue.
	go h.persistRequestHistory(ctx, persistingDone, responseChanToPersist, errChan)

	responses := []*model.Activity{}
	/* loop until either of the following gets completed:
	* - time expires
	* - an error is encountered
	* - persisting process is completed.
	 */
	for {
		select {
		case <-time.After(2 * time.Second):
			return nil, errActivityApiNotAvailabe
		case err := <-errChan:
			return nil, err
		case responses = <-responseChan:
			responseChanToPersist <- responses
		case <-persistingDone:
			return responses, nil
		}
	}

}

func (h *httpGetEventsHandler) prepareResponse(ctx context.Context, res chan []*model.Activity, errChan chan error) {
	respChan := make(chan *responseMeta)
	var wg sync.WaitGroup
	done := make(chan struct{})

	// call the api three times
	for i := 0; i < 3; i++ {
		wg.Add(1)

		go h.callActivityAPI(ctx, &wg, respChan)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	responses := []*model.Activity{}
	/* loop until either of the following gets completed:
	* - time expires
	* - an error is encountered
	* - all the responses are collected, which is marked by the done chan.
	 */
	for {
		select {
		case res := <-respChan:
			if res.err != nil {
				err := errors.Wrap(res.err, "call api")
				errChan <- err
				return
			}
			responses = append(responses, res.response)
		case <-done:
			res <- responses
			return
		case <-ctx.Done():
			return
		}
	}
}

func (h *httpGetEventsHandler) callActivityAPI(ctx context.Context, wg *sync.WaitGroup, resp chan *responseMeta) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
		res, err := h.fetchActivity(ctx)
		resp <- &responseMeta{response: res, err: err}
		return
	}
}

func (h *httpGetEventsHandler) persistRequestHistory(ctx context.Context, done chan struct{}, respChan chan []*model.Activity, errChan chan error) {
	select {
	case <-ctx.Done():
		return
	case res := <-respChan:
		err := h.persistResponse(ctx, res)
		if err != nil {
			err = errors.Wrap(err, "persist response")
			errChan <- err
			return
		}
		done <- struct{}{}
	}
}
