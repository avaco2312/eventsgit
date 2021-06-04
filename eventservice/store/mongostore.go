package store

import (
	"eventsgit/contracts"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoStore struct {
	session *mgo.Session
	db      string
}

func NewMongoStore(conn string, ndb string) (*mongoStore, error) {
	ses, err := mgo.Dial(conn)
	if err != nil {
		return nil, err
	}
	return &mongoStore{session: ses, db: ndb}, nil
}

func (m *mongoStore) SearchId(id interface{}) (*contracts.Event, error) {
	ses := m.session.Copy()
	defer ses.Close()
	var event contracts.Event
	err := ses.DB(m.db).C("events").FindId(bson.ObjectId(id.([]byte))).One(&event)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (m *mongoStore) SearchName(nEvent string) (*contracts.Event, error) {
	ses := m.session.Copy()
	defer ses.Close()
	var event contracts.Event
	err := ses.DB(m.db).C("events").Find(bson.M{"name": nEvent}).One(&event)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (m *mongoStore) SearchAll() ([]contracts.Event, error) {
	ses := m.session.Copy()
	defer ses.Close()
	var events []contracts.Event
	err := ses.DB(m.db).C("events").Find(nil).All(&events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (m *mongoStore) AddEvent(ev *contracts.Event) (interface{}, error) {
	ses := m.session.Copy()
	defer ses.Close()
	ev.ID = bson.NewObjectId()
	ev.Location.ID = bson.NewObjectId()
	err := ses.DB(m.db).C("events").Insert(ev)
	if err != nil {
		return nil, err
	}
	return string(ev.ID.(bson.ObjectId)), nil
}
