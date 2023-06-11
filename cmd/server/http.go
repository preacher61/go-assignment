package main

import (
	"log"
	"net/http"
	"preacher61/go-assignment/httpjson"
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
		name: "get-activities",
		handler: func() http.Handler {
			return newHTTPGetEventsHandler()
		},
		configure: func(r *mux.Route) {
			r.Methods(http.MethodGet).Path("/activities")
		},
	},
}

func httpNotFound(w http.ResponseWriter, req *http.Request) {
	httpjson.WriteResponse(w, http.StatusNotFound, &errorResponse{
		Error: "HTTP route not found",
	})
}

func httpMethodNotAllowed(w http.ResponseWriter, req *http.Request) {
	httpjson.WriteResponse(w, http.StatusMethodNotAllowed, &errorResponse{
		Error: "HTTP method not allowed",
	})
}

func panicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				log.Printf("recovering from err %v\n %s", err, buf)
				httpjson.WriteResponse(w, http.StatusInternalServerError, &errorResponse{
					Error: "our server got panic",
				})
			}
		}()

		h.ServeHTTP(w, r)
	})
}

type errorResponse struct {
	Error string `json:"error"`
}
