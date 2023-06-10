package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
)

func newHTTPRouter() *mux.Router {
	r := mux.NewRouter()
	registerMiddleWares(r)
	registerRoutes(r)
	return r
}

func registerMiddleWares(r *mux.Router) {
	r.Use(panicRecovery)
}

func registerRoutes(rr *mux.Router) {
	rr.NotFoundHandler = http.HandlerFunc(httpNotFound)
	rr.MethodNotAllowedHandler = http.HandlerFunc(httpMethodNotAllowed)
	for _, val := range routes {
		handler := val.handler()

		r := rr.NewRoute()
		r.Name(val.name)
		r.Handler(handler)
		val.configure(r)
	}
}

var routes = []struct {
	name      string
	handler   func() http.Handler
	configure func(r *mux.Route)
}{
	{
		name: "get-events",
		handler: func() http.Handler {
			return newHTTPGetEventsHandler()
		},
		configure: func(r *mux.Route) {
			r.Methods(http.MethodGet).Path("/events")
		},
	},
}

func httpNotFound(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "HTTP route not found"}`))
}

func httpMethodNotAllowed(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(`{"error": "HTTP method not allowed"}`))
}

func panicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				log.Printf("recovering from err %v\n %s", err, buf)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"error":"our server got panic"}`))
			}
		}()

		h.ServeHTTP(w, r)
	})
}
