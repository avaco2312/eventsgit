package main

import (
	"flag"
	"log"

	"eventsgit/aws"
	"eventsgit/bookservice/rest"
	"eventsgit/bookservice/store"
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
	processor, err := msgqueue.NewEventProcessor(sConf.MQueueType, driver, sConf.MQueueExchange,
		sConf.MQueueQueue, &rest.ListenHandler{Store: store}, rest.StaticEventMapper{}, "event.Created")
	if err != nil {
		log.Fatalf("error: imposible conectar MQueue: %v", err)
	}
	go func() {
		processor.ProcessEvents()
	}()
	cherr := rest.ServeApi(store, sConf.RestfulEndpoint, sConf.EndpointPath)
	<-cherr
}
