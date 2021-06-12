package rest

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"eventsgit/bookservice/store"
)

func ServeApi(store store.Store, endpoint string, path string) chan error {
	go func() {
		h := http.NewServeMux()
		h.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9100", handlers.CORS()(h))
	}()
	handler, _ := NewHandler(store)
	r := mux.NewRouter()
	eventsRouter := r.PathPrefix("/" + path).Subrouter()
	eventsRouter.Methods("GET").Path("/events/{search}/{params}").HandlerFunc(handler.search)
	eventsRouter.Methods("GET").Path("/events").HandlerFunc(handler.searchAll)
	cherr := make(chan error)
	go func() {
		cherr <- http.ListenAndServe(endpoint, handlers.CORS()(r))
		//cherr <- http.ListenAndServeTLS(endpoint, "cert.pem", "key.pem", r)
	}()
	return cherr
}
