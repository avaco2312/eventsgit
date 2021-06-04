package store

import (
	awsses "eventsgit/aws"
	"eventsgit/contracts"
	"encoding/hex"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type eventAWS struct {
	ID   string
	Name string `dynamodbav:"EventName"`
	Body string
}

type DynamoStore struct {
	service *dynamodb.DynamoDB
	db      string
}

func NewDynamoStore(db string) (*DynamoStore, error) {
	return &DynamoStore{
		service: dynamodb.New(awsses.Sesion),
		db:      db,
	}, nil
}

func (d *DynamoStore) AddEvent(event *contracts.Event) error {
	evaws := eventAWS{}
	evaws.ID = event.ID.(string)
	evaws.Name = event.Name
	sbody, err := json.Marshal(event)
	if err != nil {
		return err
	}
	evaws.Body = string(sbody)
	av, err := dynamodbattribute.MarshalMap(evaws)
	if err != nil {
		return err
	}
	_, err = d.service.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.db),
		Item:      av,
	})
	return err
}

func (d *DynamoStore) SearchId(id interface{}) (*contracts.Event, error) {
	bid := hex.EncodeToString([]byte(id.([]byte)))
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(bid),
			},
		},
		TableName: aws.String(d.db),
	}
	result, err := d.service.GetItem(input)
	if err != nil {
		return nil, err
	}
	evaws := eventAWS{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &evaws)
	if err != nil {
		return nil, err
	}
	var event contracts.Event
	err = json.Unmarshal([]byte(evaws.Body), &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (d *DynamoStore) SearchName(name string) (*contracts.Event, error) {
	input := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("EventName = :n"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {
				S: aws.String(name),
			},
		},
		IndexName: aws.String("EventName-index"),
		TableName: aws.String(d.db),
	}
	result, err := d.service.Query(input)
	if err != nil {
		return nil, err
	}
	evaws := eventAWS{}
	var event contracts.Event
	if len(result.Items) > 0 {
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &evaws)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(evaws.Body), &event)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return &event, err
}

func (d *DynamoStore) SearchAll() ([]contracts.Event, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(d.db),
	}
	result, err := d.service.Scan(input)
	if err != nil {
		return nil, err
	}
	var evawss []eventAWS
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &evawss)
	if err != nil {
		return nil, err
	}
	var events []contracts.Event
	var event contracts.Event
	for _, evaws := range evawss {
		err = json.Unmarshal([]byte(evaws.Body), &event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, err
}
