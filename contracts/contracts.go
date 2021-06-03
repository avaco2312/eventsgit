package contracts

type Event struct {
	ID        interface{} `bson:"_id"`
	Name      string `dynamodbav:"EventName"`
	Duration  int
	StartDate int64
	EndDate   int64
	Location  Location
}

type Location struct {
	ID        interface{} `bson:"_id"`
	Name      string `dynamodbav:"LocationName"`
	Address   string
	Country   string
	OpenTime  int
	CloseTime int
	Halls     []Hall
}

type Hall struct {
	Name     string `dynamodbav:"HallName"`
	Location string
	Capacity int
}

type EventCreated struct {
	Event *Event
}

// EventName returns the event's name
func (e EventCreated) EventName() string {
	return "event.Created"
}