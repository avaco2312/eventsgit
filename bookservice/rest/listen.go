package rest

import (
	"encoding/json"
	"fmt"
	"log"

	"Cloud-libro/Chapter04/yo/bookservice/store"
	"Cloud-libro/Chapter04/yo/contracts"
	"Cloud-libro/Chapter04/yo/msgqueue"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2/bson"
)

type ListenHandler struct {
	Store store.Store
}

func (l *ListenHandler) Handle(event interface{}) {
	switch e := event.(type) {
	case *contracts.EventCreated:
		if !bson.IsObjectIdHex(e.Event.ID.(string)) {
			log.Printf("event %v did not contain valid object ID", e)
			return
		}
		l.Store.AddEvent(e.Event)
	default:
		log.Printf("unknown event type: %T", e)
	}
}

type StaticEventMapper struct{}

func (e StaticEventMapper) MapEvent(eventName string, serialized interface{}) (msgqueue.Event, error) {
	var event *contracts.EventCreated
	switch eventName {
	case "event.Created":
		event = &contracts.EventCreated{}
	default:
		return nil, fmt.Errorf("unknown event type %s", eventName)
	}
	switch s := serialized.(type) {
	case []byte:
		err := json.Unmarshal(s, event)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
		}
	default:
		cfg := mapstructure.DecoderConfig{
			Result:  event,
			TagName: "json",
		}
		dec, err := mapstructure.NewDecoder(&cfg)
		if err != nil {
			return nil, fmt.Errorf("could not initialize decoder for event %s: %s", eventName, err)
		}
		err = dec.Decode(s)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal event %s: %s", eventName, err)
		}
	}
	return event, nil
}
