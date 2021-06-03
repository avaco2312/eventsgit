package store

import (
	"Cloud-libro/Chapter04/yo/aws"
	"Cloud-libro/Chapter04/yo/contracts"
	"fmt"
)

type Store interface {
	SearchId(interface{}) (*contracts.Event, error)
	SearchName(string) (*contracts.Event, error)
	SearchAll() ([]contracts.Event, error)
	AddEvent(*contracts.Event) (interface{}, error)
}

func NewStore(dbType string, connString string, db string) (Store, error) {
	var st Store
	var err error
	switch dbType {
	case "mongo":
		st, err = NewMongoStore(connString, db)
	case "dynamo":
		err = aws.SetSession()
		if err == nil {
			st, err = NewDynamoStore(db)
		}
	default:
		return nil, fmt.Errorf("error: Unknown DB driver %s", dbType)
	}
	if err != nil {
		return nil, err
	}
	return st, nil
}
