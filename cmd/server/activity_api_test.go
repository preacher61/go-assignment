package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestActivityAPISuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := `{
			"activity": "test activity",
    		"key": "4894697"
		}`
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(resp))
	}))
	defer server.Close()

	a := &activityAPI{
		url:         server.URL,
		rateLimiter: rate.NewLimiter(rate.Every(1*time.Second), 1),
	}

	res, err := a.getEvent(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if res.Key != "4894697" {
		t.Fatalf("invalid response.key expected:4894697 got: %s", res.Key)
	}
}

func TestActivityAPIInternalServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	a := &activityAPI{
		url:         server.URL,
		rateLimiter: rate.NewLimiter(rate.Every(1*time.Second), 1),
	}

	_, err := a.getEvent(context.TODO())
	if err == nil {
		t.Fatal("error expected")
	}
}

func TestActivityAPISuccessUnMarshalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := `{
			"activity": "test activity",
    		"key": 4894697
		}`
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(resp))
	}))
	defer server.Close()

	a := &activityAPI{
		url:         server.URL,
		rateLimiter: rate.NewLimiter(rate.Every(1*time.Second), 1),
	}

	_, err := a.getEvent(context.TODO())
	if err == nil {
		t.Fatal("error expected")
	}
}

func TestActivityAPISuccessNilResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("null"))
	}))
	defer server.Close()

	a := &activityAPI{
		url:         server.URL,
		rateLimiter: rate.NewLimiter(rate.Every(1*time.Second), 1),
	}

	_, err := a.getEvent(context.TODO())
	if err == nil {
		t.Fatal("error expected")
	}
}
