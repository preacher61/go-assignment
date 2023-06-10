package main

import "net/http"

type httpGetEventsHandler struct{}

func newHTTPGetEventsHandler() *httpGetEventsHandler {
	return &httpGetEventsHandler{}
}

func (h *httpGetEventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("panic occured")
}
