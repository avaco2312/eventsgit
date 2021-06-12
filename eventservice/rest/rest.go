package rest

import (
	"net/http"

	"eventsgit/eventservice/store"
	"eventsgit/msgqueue"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ServeApi(store store.Store, emitter msgqueue.EventEmitter, endpoint string, path string) chan error {
	go func() {
		h := http.NewServeMux()
		h.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9100", handlers.CORS()(h))
	}()
	handler, _ := NewHandler(store, emitter)
	r := mux.NewRouter()
	eventsRouter := r.PathPrefix("/" + path).Subrouter()
	eventsRouter.Methods("GET").Path("/{search}/{params}").HandlerFunc(handler.search)
	eventsRouter.Methods("GET").Path("").HandlerFunc(handler.searchAll)
	eventsRouter.Methods("POST").Path("").HandlerFunc(handler.addEvent)
	cherr := make(chan error)
	go func() {
		cherr <- http.ListenAndServe(endpoint, handlers.CORS()(r))
		//cherr <- http.ListenAndServeTLS(endpoint, "cert.pem", "key.pem", r)
	}()
	return cherr
}
