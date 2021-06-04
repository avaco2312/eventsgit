package main

import (
	"flag"
	"log"

	"eventsgit/aws"
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
	store, err := store.NewStore(sConf.Databasetype, sConf.DBConnection, sConf.DbName)
	if err != nil {
		log.Fatalf("error: imposible conectar la BD: %v", err)
	}
	var driver string
	switch sConf.MQueueType {
	case "amqp":
		driver = sConf.AMQPMessageBroker
	case "kafka":
		driver = sConf.KafkaMessageBroker
	case "sqs":
		driver = ""
		err = aws.SetSession()
		if err != nil {
			log.Fatalf("error: imposible conectar AWS: %v", err)
		}
	default:
		log.Fatalf("error: MQueue driver desconocido %s", sConf.MQueueType)
	}
	emitter, err := msgqueue.NewEventEmitter(sConf.MQueueType, driver, sConf.MQueueExchange)
	if err != nil {
		log.Fatalf("error: imposible conectar MQueue: %v", err)
	}
	cherr := rest.ServeApi(store, emitter, sConf.RestfulEndpoint, sConf.EndpointPath)
	<-cherr
}
