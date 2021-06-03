package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	DBTypeDefault              = "mongo" // "mongo" or "dynamo"
	DBConnectionDefault        = "root:example@127.0.0.1"
	RestfulEndpointDefault     = ":8070"
	EndpointPathDefault        = "events"
	DbNameDefault              = "myevents"
	MQueueTypeDefault          = "kafka" // "amqp" or "kafka"
	AMQPMessageBrokerDefault   = "amqp://127.0.0.1:5672"
	KafkaMessageBrokersDefault = "localhost:9092"
	MQueueExchangeDefault      = "events"
	EnvDefault                 = "no"
)

type ServiceConfig struct {
	Databasetype       string `json:"databasetype"`
	DBConnection       string `json:"dbconnection"`
	RestfulEndpoint    string `json:"restfulapi_endpoint"`
	EndpointPath       string `json:"endpoint_path"`
	DbName             string `json:"dbname"`
	MQueueType         string `json:"mqueuetype"`
	AMQPMessageBroker  string `json:"amqpmessagebroker"`
	KafkaMessageBroker string `json:"kafkamessagebroker"`
	MQueueExchange     string `json:"mqueueexchange"`
	Env                string `json:"env"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEndpointDefault,
		EndpointPathDefault,
		DbNameDefault,
		MQueueTypeDefault,
		AMQPMessageBrokerDefault,
		KafkaMessageBrokersDefault,
		MQueueExchangeDefault,
		EnvDefault,
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Configuration file not found. Continuing with default values.")
		return conf, nil
	}
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
			return conf, err
	}
	if conf.Env != "no" {
		conf.DBConnection = "root:example@" + os.Getenv("MONGO_SERVICE_HOST")
		switch conf.MQueueType {
		case "amqp":
			conf.AMQPMessageBroker="amqp://"+os.Getenv("AMQP_SERVICE_HOST")+":5672"
		case "kafka":
			conf.KafkaMessageBroker=os.Getenv("KAFKA_SERVICE_HOST")+":9092"
		default:
			return conf, fmt.Errorf("error: Unknown MQueue type %s", conf.MQueueType)
		}
	}
	return conf, nil
}
