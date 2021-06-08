package main

import (
	"encoding/json"
	"eventsgit/aws"
	"fmt"
	"log"
	"os"
)

const (
	DBTypeDefault          = "mongo" // "mongo" or "dynamo"
	DBConnectionDefault    = ""
	RestfulEndpointDefault = ":8070"
	EndpointPathDefault    = "events"
	DbNameDefault          = "myevents"
	MQueueTypeDefault      = "kafka" // "amqp" or "kafka" or "sqs"
	MQueueExchangeDefault  = "events"
	EnvDefault             = "local" // local, docker, kubernet
	MQueueDriverDefault    = ""
)

type ServiceConfig struct {
	DBType          string `json:"databasetype"`
	DBConnection    string `json:"dbconnection"`
	RestfulEndpoint string `json:"restfulapi_endpoint"`
	EndpointPath    string `json:"endpoint_path"`
	DbName          string `json:"dbname"`
	MQueueType      string `json:"mqueuetype"`
	MQueueExchange  string `json:"mqueueexchange"`
	Env             string `json:"env"`
	MQueueDriver    string `json:"driver"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEndpointDefault,
		EndpointPathDefault,
		DbNameDefault,
		MQueueTypeDefault,
		MQueueExchangeDefault,
		EnvDefault,
		MQueueDriverDefault,
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
	switch conf.Env {
	case "local":
		switch conf.DBType {
		case "mongo":
			conf.DBConnection = "root:example@localhost"
		case "dynamo":
			conf.DBConnection = ""
			err = aws.SetSession()
			if err != nil {
				return conf, fmt.Errorf("error: Imposible conectar AWS: %v", err)
			}
		default:
			return conf, fmt.Errorf("error: Unknown Database type %s", conf.DBType)
		}
		switch conf.MQueueType {
		case "amqp":
			conf.MQueueDriver = "amqp://localhost:5672"
		case "kafka":
			conf.MQueueDriver = "localhost:9092"
		case "sqs":
			conf.MQueueDriver = ""
			err = aws.SetSession()
			if err != nil {
				return conf, fmt.Errorf("error: Imposible conectar AWS: %v", err)
			}
		default:
			return conf, fmt.Errorf("error: Unknown MQueue type %s", conf.MQueueType)
		}
	case "docker":
		switch conf.DBType {
		case "mongo":
			conf.DBConnection = "root:example@mongo"
		case "dynamo":
			conf.DBConnection = ""
			err = aws.SetSession()
			if err != nil {
				return conf, fmt.Errorf("error: Imposible conectar AWS: %v", err)
			}
		default:
			return conf, fmt.Errorf("error: Unknown Database type %s", conf.DBType)
		}
		switch conf.MQueueType {
		case "amqp":
			conf.MQueueDriver = "amqp://rabitmq:5672"
		case "kafka":
			conf.MQueueDriver = "kafka:9092"
		case "sqs":
			conf.MQueueDriver = ""
			err = aws.SetSession()
			if err != nil {
				return conf, fmt.Errorf("error: Imposible conectar AWS: %v", err)
			}
		default:
			return conf, fmt.Errorf("error: Unknown MQueue type %s", conf.MQueueType)
		}
	case "kubernet":
		switch conf.DBType {
		case "mongo":
			conf.DBConnection = "root:example@" + os.Getenv("MONGO_SERVICE_HOST")
		case "dynamo":
			conf.DBConnection = ""
			err = aws.SetSession()
			if err != nil {
				return conf, fmt.Errorf("error: Imposible conectar AWS: %v", err)
			}
		}
		switch conf.MQueueType {
		case "amqp":
			conf.MQueueDriver = "amqp://" + os.Getenv("AMQP_SERVICE_HOST") + ":5672"
		case "kafka":
			conf.MQueueDriver = os.Getenv("KAFKA_SERVICE_HOST") + ":9092"
		case "sqs":
			conf.MQueueDriver = ""
			err = aws.SetSession()
			if err != nil {
				return conf, fmt.Errorf("error: Imposible conectar AWS: %v", err)
			}
		default:
			return conf, fmt.Errorf("error: Unknown MQueue type %s", conf.MQueueType)
		}
	}
	return conf, nil
}
