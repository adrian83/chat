package db

import (
	"strconv"

	r "github.com/dancannon/gorethink"
)

// Config contains database configuration properties.
type Config struct {
	Host             string
	Port             int
	DBName           string
	UsersTableName   string
	UsersTablePKName string
}

// DataBase struct used for working with RethinkDB
type DataBase struct {
	Config  *Config
	Session *r.Session
}

// New returns pointer to new instance of DataBase and error
func New(config *Config) (*DataBase, error) {
	session, err := r.Connect(r.ConnectOpts{
		Address:  config.Host + ":" + strconv.Itoa(config.Port),
		Database: config.DBName,
	})
	if err != nil {
		return nil, err
	}

	return &DataBase{Session: session, Config: config}, nil
}

// Close closes RethinkDB session
func (db *DataBase) Close() error {
	return db.Session.Close()
}

func (db *DataBase) createUUID() (string, error) {
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
func (db *DataBase) Setup() error {
	if err := db.createDbIfItDoesNotExist(db.Config.DBName); err != nil {
		return err
	}
	if err := db.createTableIfItDoesNotExist(db.Config.DBName, db.Config.UsersTableName, db.Config.UsersTablePKName); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) createDbIfItDoesNotExist(dbName string) error {
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

func (db *DataBase) createTableIfItDoesNotExist(dbName, tableName, pkName string) error {

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

func (db *DataBase) boolResp(cursor *r.Cursor) (bool, error) {
	resp := new(bool)
	if err := cursor.One(resp); err != nil {
		return false, err
	}
	return *resp, nil
}
