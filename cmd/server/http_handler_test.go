package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"preacher61/go-assignment/model"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestHTTPHanlderSuccess(t *testing.T) {
	h := &httpGetEventsHandler{
		fetchActivity: func(ctx context.Context) (*model.Activity, error) {
			return &model.Activity{
				Activity: "test activity",
				Key:      "6786876",
			}, nil
		},
		persistResponse: func(ctx context.Context, response []*model.Activity) error {
			return nil
		},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost/activities", http.NoBody)
	h.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusOK)
	}

	b, err := io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	var res []*model.Activity
	err = json.Unmarshal(b, &res)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 3 {
		t.Fatalf("invalid response expected length:3 got: %d", len(res))
	}
}

func TestHTTPHanlderErrorFromOneRequest(t *testing.T) {
	var callerCounter atomic.Int32
	h := &httpGetEventsHandler{
		fetchActivity: func(ctx context.Context) (*model.Activity, error) {
			if callerCounter.Load() == 1 {
				return nil, errors.New("some i/o error")
			}
			callerCounter.Add(1)
			return &model.Activity{
				Activity: "test activity",
				Key:      "6786876",
			}, nil
		},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost/activities", http.NoBody)
	h.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusInternalServerError)
	}

	b, err := io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	var res errorResponse
	err = json.Unmarshal(b, &res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Error == "" {
		t.Fatal("error expected")
	}
}

func TestHTTPHanlderFetchActivityTimeOutErr(t *testing.T) {
	h := &httpGetEventsHandler{
		fetchActivity: func(ctx context.Context) (*model.Activity, error) {
			time.Sleep(1 * time.Minute)
			return &model.Activity{
				Activity: "test activity",
				Key:      "6786876",
			}, nil
		},
		persistResponse: func(ctx context.Context, response []*model.Activity) error {
			return nil
		},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost/activities", http.NoBody)
	h.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusInternalServerError)
	}

	b, err := io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	var res errorResponse
	err = json.Unmarshal(b, &res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Error == "" {
		t.Fatal("error expected")
	}
}

func TestHTTPHanlderPersistResponseTimeOutErr(t *testing.T) {
	h := &httpGetEventsHandler{
		fetchActivity: func(ctx context.Context) (*model.Activity, error) {
			return &model.Activity{
				Activity: "test activity",
				Key:      "6786876",
			}, nil
		},
		persistResponse: func(ctx context.Context, response []*model.Activity) error {
			time.Sleep(1 * time.Minute)
			return nil
		},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost/activities", http.NoBody)
	h.ServeHTTP(w, req)
	w.Flush()
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected code: got %d, want %d", w.Code, http.StatusInternalServerError)
	}

	b, err := io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	var res errorResponse
	err = json.Unmarshal(b, &res)
	if err != nil {
		t.Fatal(err)
	}
	if res.Error == "" {
		t.Fatal("error expected")
	}
}
