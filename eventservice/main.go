package main

import (
	"flag"
	"log"

	"eventsgit/eventservice/rest"
	"eventsgit/eventservice/store"
	"eventsgit/msgqueue"
)

func main() {
	confPath := flag.String("conf", "config.json", "flag to set the configuration json file")
	flag.Parse()
	sConf, err := ExtractConfiguration(*confPath)
	if err != nil {
		log.Fatalf("error: imposible cargar configuración: %v", err)
	}
	store, err := store.NewStore(sConf.DBType, sConf.DBConnection, sConf.DbName)
	if err != nil {
		log.Fatalf("error: imposible conectar la BD: %v", err)
	}
	emitter, err := msgqueue.NewEventEmitter(sConf.MQueueType, sConf.MQueueDriver, sConf.MQueueExchange)
	if err != nil {
		log.Fatalf("error: imposible conectar MQueue: %v", err)
	}
	cherr := rest.ServeApi(store, emitter, sConf.RestfulEndpoint, sConf.EndpointPath)
	<-cherr
}
