package main

import (
	"log"
	"net/http"

	"eventsgit/bookservice/rest"
	"eventsgit/bookservice/store"
	"eventsgit/msgqueue"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	sConf, err := ExtractConfiguration()
	if err != nil {
		log.Fatalf("error: imposible cargar configuración: %v", err)
	}
	store, err := store.NewStore(sConf.dbType, sConf.dbConnection, sConf.dbName)
	if err != nil {
		log.Fatalf("error: imposible conectar la BD: %v", err)
	}
	processor, err := msgqueue.NewEventProcessor(sConf.queueType, sConf.queueDriver, sConf.queueExchange,
		sConf.queueQueue, &rest.ListenHandler{Store: store}, rest.StaticEventMapper{}, "event.Created")
	if err != nil {
		log.Fatalf("error: imposible conectar MQueue: %v", err)
	}
	go func() {
		processor.ProcessEvents()
	}()
	go func() {
		h := http.NewServeMux()
		h.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9100", h)
	}()
	cherr := rest.ServeApi(store, sConf.restfulEndpoint, sConf.endpointPath)
	<-cherr
}
