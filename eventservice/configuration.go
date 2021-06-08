package main

import (
	"eventsgit/aws"
	"fmt"
	"os"
)

const (
	dbTypeDefault          = "mongo" // "mongo" or "dynamo"
	dbConnectionDefault    = ""
	restfulEndpointDefault = ":8070"
	endpointPathDefault    = "events"
	dbNameDefault          = "myevents"
	queueTypeDefault       = "kafka" // "amqp" or "kafka" or "sqs"
	queueExchangeDefault   = "events"
	envDefault             = "local" // local, docker, kubernet
	queueDriverDefault     = ""
)

type ServiceConfig struct {
	dbType          string 
	dbConnection    string 
	RestfulEndpoint string 
	endpointPath    string 
	dbName          string 
	queueType       string 
	queueExchange   string 
	Env             string 
	queueDriver     string 
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		dbTypeDefault,
		dbConnectionDefault,
		restfulEndpointDefault,
		endpointPathDefault,
		dbNameDefault,
		queueTypeDefault,
		queueExchangeDefault,
		envDefault,
		queueDriverDefault,
	}
	if s, ok := os.LookupEnv("DB_TYPE"); ok {
		conf.dbType = s
	}
	if s, ok := os.LookupEnv("REST_ENDPOINT"); ok {
		conf.endpointPath = s
	}
	if s, ok := os.LookupEnv("ENDPOINT_PATH"); ok {
		conf.endpointPath = s
	}
	if s, ok := os.LookupEnv("DB_NAME"); ok {
		conf.dbName = s
	}
	if s, ok := os.LookupEnv("QUEUE_TYPE"); ok {
		conf.queueType = s
	}
	if s, ok := os.LookupEnv("QUEUE_EXCHANGE"); ok {
		conf.queueExchange = s
	}
	if s, ok := os.LookupEnv("RUN_ENV"); ok {
		conf.Env = s
	}
	switch conf.Env {
	case "local":
		switch conf.dbType {
		case "mongo":
			conf.dbConnection = "root:example@localhost"
		case "dynamo":
			conf.dbConnection = ""
			err := aws.SetSession()
			if err != nil {
				return conf, fmt.Errorf("error: Imposible conectar AWS: %v", err)
			}
		default:
			return conf, fmt.Errorf("error: Unknown Database type %s", conf.dbType)
		}
		switch conf.queueType {
		case "amqp":
			conf.queueDriver = "amqp://localhost:5672"
		case "kafka":
			conf.queueDriver = "localhost:9092"
		case "sqs":
			conf.queueDriver = ""
			err := aws.SetSession()
			if err != nil {
				return conf, fmt.Errorf("error: Imposible conectar AWS: %v", err)
			}
		default:
			return conf, fmt.Errorf("error: Unknown MQueue type %s", conf.queueType)
		}
	case "docker":
		switch conf.dbType {
		case "mongo":
			conf.dbConnection = "root:example@mongo"
		case "dynamo":
			conf.dbConnection = ""
			err := aws.SetSession()
			if err != nil {
				return conf, fmt.Errorf("error: Imposible conectar AWS: %v", err)
			}
		default:
			return conf, fmt.Errorf("error: Unknown Database type %s", conf.dbType)
		}
		switch conf.queueType {
		case "amqp":
			conf.queueDriver = "amqp://rabitmq:5672"
		case "kafka":
			conf.queueDriver = "kafka:9092"
		case "sqs":
			conf.queueDriver = ""
			err := aws.SetSession()
			if err != nil {
				return conf, fmt.Errorf("error: Imposible conectar AWS: %v", err)
			}
		default:
			return conf, fmt.Errorf("error: Unknown MQueue type %s", conf.queueType)
		}
	case "kubernet":
		switch conf.dbType {
		case "mongo":
			conf.dbConnection = "root:example@" + os.Getenv("MONGO_SERVICE_HOST")
		case "dynamo":
			conf.dbConnection = ""
			err := aws.SetSession()
			if err != nil {
				return conf, fmt.Errorf("error: Imposible conectar AWS: %v", err)
			}
		}
		switch conf.queueType {
		case "amqp":
			conf.queueDriver = "amqp://" + os.Getenv("AMQP_SERVICE_HOST") + ":5672"
		case "kafka":
			conf.queueDriver = os.Getenv("KAFKA_SERVICE_HOST") + ":9092"
		case "sqs":
			conf.queueDriver = ""
			err := aws.SetSession()
			if err != nil {
				return conf, fmt.Errorf("error: Imposible conectar AWS: %v", err)
			}
		default:
			return conf, fmt.Errorf("error: Unknown MQueue type %s", conf.queueType)
		}
	}
	return conf, nil
}
