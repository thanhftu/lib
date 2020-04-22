package mongolayer

import (
	"github.com/thanhftu/lib/persistence"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB     = "myevents"
	USERS  = "users"
	EVENTS = "events"
)

type MongoDBlayer struct {
	session *mgo.Session
}

func NewMongoDBLayer(connection string) (*MongoDBlayer, error) {
	s, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}
	return &MongoDBlayer{
		session: s,
	}, err
}

func (mgoLayer *MongoDBlayer) AddEvent(e persistence.Event) ([]byte, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()

	if !e.ID.Valid() {
		e.ID = bson.NewObjectId()
	}
	if !e.Location.ID.Valid() {
		e.Location.ID = bson.NewObjectId()
	}
	return []byte(e.ID), s.DB(DB).C(EVENTS).Insert(e)
}

func (mgoLayer *MongoDBlayer) FindEvent(id []byte) (persistence.Event, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()

	e := persistence.Event{}                                    //create empty event object
	err := s.DB(DB).C(EVENTS).FindId(bson.ObjectId(id)).One(&e) //e will obtain data of the retrieved document
	return e, err
}

func (mgoLayer *MongoDBlayer) FindEventByName(name string) (persistence.Event, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()

	e := persistence.Event{}
	err := s.DB(DB).C(EVENTS).Find(bson.M{"name": name}).One(&e) //M standfor Map
	return e, err
}

func (mgoLayer *MongoDBlayer) FindAllAvailableEvents() ([]persistence.Event, error) {
	s := mgoLayer.getFreshSession()
	defer s.Close()

	events := []persistence.Event{}
	err := s.DB(DB).C(EVENTS).Find(nil).All(&events) // Find(nil) return every thing found
	return events, err
}

// getFreshSession helper method implemented in our code to
// help retrieve a fresh database session from the connection pool

func (mgoLayer *MongoDBlayer) getFreshSession() *mgo.Session {
	return mgoLayer.session.Copy()
}
