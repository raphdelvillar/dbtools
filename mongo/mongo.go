package mongo

import (
	"dbtools/postgres"
	"fmt"

	"github.com/globalsign/mgo"
)

var pgdb postgres.Postgres
var mdb Mongo

// Mongo --
type Mongo struct {
	URL        string
	Session    *mgo.Session
	Database   string
	Collection string
}

// Init --
func (m *Mongo) Init() {
	username := ""
	password := ""
	host := ""
	port := ""
	mdb.URL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)

	mdb.Connect()
}

// Connect --
func (m *Mongo) Connect() error {
	session, err := mgo.Dial(m.URL)

	if err != nil {
		return err
	}

	m.Session = session
	return nil
}

// Disconnect --
func (m *Mongo) Disconnect() {
	m.Session.Close()
}

// Insert --
func (m *Mongo) Insert(data interface{}) error {
	c := m.Session.DB(m.Database).C(m.Collection)
	err := c.Insert(data)

	if err != nil {
		return err
	}

	return nil
}
