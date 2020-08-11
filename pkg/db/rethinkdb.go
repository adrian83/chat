package db

import (
	"fmt"
	"strconv"

	logger "github.com/sirupsen/logrus"
	r "gopkg.in/gorethink/gorethink.v4"
)

const (
	usersTableName    = "users"
	usersTableNameKey = "name"
)

// RethinkDB is a struct that allows communication with RethinkDB.
type RethinkDB struct {
	host string
	name string
	port int

	session *r.Session
}

// NewRethinkDB returns new instance of RethinkDB (not connected).
func NewRethinkDB(host string, port int, name string) *RethinkDB {
	return &RethinkDB{
		host: host,
		name: name,
		port: port,
	}
}

// Connect creates connection to RethinkDB.
func (rt *RethinkDB) Connect() error {
	session, err := r.Connect(r.ConnectOpts{
		Address:  rt.host + ":" + strconv.Itoa(rt.port),
		Database: rt.name,
	})
	if err != nil {
		return fmt.Errorf("cannot connect to RethingDB, error: %w", err)
	}

	rt.session = session
	return nil
}

// Setup creates database and tables if they don't exist.
func (rt *RethinkDB) Setup() error {
	dbExists, err := rt.containsDB()
	if err != nil {
		return fmt.Errorf("cannot check if database exist, error: %w", err)
	}

	if !dbExists {
		if err := rt.createDB(); err != nil {
			return fmt.Errorf("cannot create database, error: %w", err)
		}
	}

	tableExists, err := rt.containsTable(usersTableName)
	if err != nil {
		return fmt.Errorf("cannot check if table %v exist, error: %w", usersTableName, err)
	}

	if !tableExists {
		if err := rt.createTable(usersTableName, usersTableNameKey); err != nil {
			return err
		}
	}

	return nil
}

// Close closes connection to RethinkDB.
func (rt *RethinkDB) Close() {
	if err := rt.session.Close(); err != nil {
		logger.Errorf("cannot close RethinkDB session, error: %v", err)
	}
}

// UUID returns new UUID.
func (rt *RethinkDB) UUID() (string, error) {
	cursor, err := r.UUID().Run(rt.session)
	if err != nil {
		return "", fmt.Errorf("cannot create UUID, error: %w", err)
	}

	var uuid string
	if err := cursor.One(&uuid); err != nil {
		return "", fmt.Errorf("cannot fetch UUID from cursor, error: %w", err)
	}

	return uuid, nil
}

func (rt *RethinkDB) createTable(tableName, primaryKey string) error {
	_, err := r.DB(rt.name).TableCreate(tableName, r.TableCreateOpts{PrimaryKey: primaryKey}).Run(rt.session)

	return fmt.Errorf("cannot create table %v with primary key %v, error: %w", tableName, primaryKey, err)
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
	var resp bool
	if err := cursor.One(&resp); err != nil {
		return false, err
	}
	return resp, nil
}

// GetUserTable returns users table.
func (rt *RethinkDB) GetUserTable() *RethinkTable {
	return &RethinkTable{
		name:    usersTableName,
		term:    r.DB(rt.name).Table(usersTableName),
		rethink: rt,
	}
}

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
func (t *RethinkTable) Find(property string, value, result interface{}) error {
	cursor, err := t.term.Filter(r.Row.Field(property).Eq(value)).Run(t.rethink.session)
	if err != nil {
		return err
	}

	if cursor.IsNil() {
		return nil
	}

	if err := cursor.One(result); err != nil {
		return err
	}

	return nil
}
