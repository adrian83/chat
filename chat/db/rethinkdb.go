package db

import (
	"fmt"
	"strconv"

	r "github.com/dancannon/gorethink"
	logger "github.com/sirupsen/logrus"
)

// RethinkDB is a struct that allows communication with RethinkDB.
type RethinkDB struct {
	host    string
	name    string
	port    int
	tables  map[string]string
	session *r.Session
}

// NewRethinkDB returns new instance of RethinkDB (not connected).
func NewRethinkDB(host string, port int, name string, tables map[string]string) *RethinkDB {
	return &RethinkDB{
		host:   host,
		name:   name,
		port:   port,
		tables: tables,
	}
}

// Connect creates connection to RethinkDB.
func (rt *RethinkDB) Connect() error {
	session, err := r.Connect(r.ConnectOpts{
		Address:  rt.host + ":" + strconv.Itoa(rt.port),
		Database: rt.name,
	})
	if err != nil {
		return err
	}

	rt.session = session
	return nil
}

// Setup creates database and tables if they don't exist.
func (rt *RethinkDB) Setup() error {
	dbExists, err := rt.containsDB()
	if err != nil {
		return err
	}

	if !dbExists {
		if err := rt.createDB(); err != nil {
			return err
		}
	}

	for tableName, primaryKey := range rt.tables {
		tableExists, err := rt.containsTable(tableName)
		if err != nil {
			return err
		}

		if !tableExists {
			if err = rt.createTable(tableName, primaryKey); err != nil {
				return err
			}
		}
	}
	return nil
}

// Close closes conection to RethinkDB.
func (rt *RethinkDB) Close() {
	if err := rt.session.Close(); err != nil {
		logger.Errorf("Error while closing RethinkDB session! Error: %v", err)
	}
}

// UUID returns new UUID.
func (rt *RethinkDB) UUID() (string, error) {
	cursor, err := r.UUID().Run(rt.session)
	if err != nil {
		return "", err
	}
	uuid := new(string)
	if err = cursor.One(uuid); err != nil {
		return "", err
	}
	return *uuid, nil
}

func (rt *RethinkDB) createTable(tableName, primaryKey string) error {
	_, err := r.DB(rt.name).TableCreate(tableName, r.TableCreateOpts{PrimaryKey: primaryKey}).Run(rt.session)
	return err
}

func (rt *RethinkDB) containsDB() (bool, error) {
	cursor, err := r.DBList().Contains(rt.name).Run(rt.session)
	if err != nil {
		return false, err
	}

	return rt.boolResp(cursor)
}

func (rt *RethinkDB) containsTable(tableName string) (bool, error) {
	cursor, err := r.DB(rt.name).TableList().Contains(tableName).Run(rt.session)
	if err != nil {
		return false, err
	}

	return rt.boolResp(cursor)
}

func (rt *RethinkDB) createDB() error {
	_, err := r.DBCreate(rt.name).Run(rt.session)
	return err
}

func (rt *RethinkDB) boolResp(cursor *r.Cursor) (bool, error) {
	resp := new(bool)
	if err := cursor.One(resp); err != nil {
		return false, err
	}
	return *resp, nil
}

// GetTable returns table with given name.
func (rt *RethinkDB) GetTable(name string) *RethinkTable {
	return &RethinkTable{
		name:    name,
		term:    r.DB(rt.name).Table(name),
		rethink: rt,
	}
}

var errNotFount = fmt.Errorf("not found")

// RethinkTable represents RethinkDB table.
type RethinkTable struct {
	name    string
	term    r.Term
	rethink *RethinkDB
}

// UUID returns new UUID.
func (t *RethinkTable) UUID() (string, error) {
	return t.rethink.UUID()
}

// Insert persists given struct into table in RethinkDB.
func (t *RethinkTable) Insert(entity interface{}) error {
	return t.term.Insert(entity).Exec(t.rethink.session)
}

// Find searches for first element with given property equal to given value.
func (t *RethinkTable) Find(property string, value interface{}, result interface{}) error {
	cursor, err := t.term.Filter(r.Row.Field(property).Eq(value)).Run(t.rethink.session)

	if cursor.IsNil() {
		return errNotFount
	}

	if err = cursor.One(result); err != nil {
		return err
	}

	return nil
}
