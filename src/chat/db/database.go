package db

import (
	"chat/config"

	"strconv"

	r "github.com/dancannon/gorethink"
)

// Database struct used for working with RethinkDB
type Database struct {
	Config  *config.DatabaseConfig
	Session *r.Session
}

// New returns pointer to new instance of Database and error
func New(config *config.DatabaseConfig) (*Database, error) {
	session, err := r.Connect(r.ConnectOpts{
		Address:  config.Host + ":" + strconv.Itoa(config.Port),
		Database: config.DBName,
	})
	if err != nil {
		return nil, err
	}

	return &Database{Session: session, Config: config}, nil
}

// Close closes RethinkDB session
func (db *Database) Close() error {
	return db.Session.Close()
}

func (db *Database) createUUID() (string, error) {
	c, err := r.UUID().Run(db.Session)
	if err != nil {
		return "", err
	}
	uuid := new(string)
	if err = c.One(uuid); err != nil {
		return "", err
	}
	return *uuid, nil
}

// Setup initializes database.
func (db *Database) Setup() error {
	if err := db.createDbIfItDoesNotExist(db.Config.DBName); err != nil {
		return err
	}
	if err := db.createTableIfItDoesNotExist(db.Config.DBName, db.Config.UsersTableName, db.Config.UsersTablePKName); err != nil {
		return err
	}
	return nil
}

func (db *Database) createDbIfItDoesNotExist(dbName string) error {
	c, err1 := r.DBList().Contains(dbName).Run(db.Session)
	if err1 != nil {
		return err1
	}
	if b, err2 := db.boolResp(c); err2 != nil {
		return err2
	} else if b {
		return nil
	}

	_, err3 := r.DBCreate(dbName).Run(db.Session)
	return err3
}

func (db *Database) createTableIfItDoesNotExist(dbName, tableName, pkName string) error {

	c, err1 := r.DB(dbName).TableList().Contains(tableName).Run(db.Session)
	if err1 != nil {
		return err1
	}

	if b, err2 := db.boolResp(c); err2 != nil {
		return err2
	} else if b {
		return nil
	}

	_, err3 := r.DB(dbName).TableCreate(tableName, r.TableCreateOpts{PrimaryKey: pkName}).Run(db.Session)
	return err3
}

func (db *Database) boolResp(cursor *r.Cursor) (bool, error) {
	resp := new(bool)
	if err := cursor.One(resp); err != nil {
		return false, err
	}
	return *resp, nil
}
