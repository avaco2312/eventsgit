package main

import (
	"fmt"
	"log"
	"net/http"

	"eventsgit/eventservice/rest"
	"eventsgit/eventservice/store"
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
		fmt.Println(err, " ", err.Error())
		log.Fatalf("error: imposible conectar la BD: %v", err)
	}
	emitter, err := msgqueue.NewEventEmitter(sConf.queueType, sConf.queueDriver, sConf.queueExchange)
	if err != nil {
		fmt.Println(err, " ", err.Error())
		log.Fatalf("error: imposible conectar MQueue: %v", err)
	}
	go func() {
		h:=http.NewServeMux()
		h.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9100", h)
	}()
	cherr := rest.ServeApi(store, emitter, sConf.restfulEndpoint, sConf.endpointPath)
	<-cherr
}
