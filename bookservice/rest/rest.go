package rest

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"../store"
)

func ServeApi(store store.Store, endpoint string, path string) chan error {

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
