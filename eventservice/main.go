package main

import (
	"fmt"
	"log"

	"eventsgit/eventservice/rest"
	"eventsgit/eventservice/store"
	"eventsgit/msgqueue"
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
	cherr := rest.ServeApi(store, emitter, sConf.restfulEndpoint, sConf.endpointPath)
	<-cherr
}
