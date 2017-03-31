package db

import (
	r "github.com/dancannon/gorethink"
)

// User is a struct containing user data.
type User struct {
	ID       string `json:"id" gorethink:"id,omitempty"`
	Name     string `json:"name" gorethink:"name,omitempty"`
	Password string `json:"password" gorethink:"password,omitempty"`
}

// Empty returns 'true' it the User struct is empty, false otherwise.
func (u *User) Empty() bool {
	return u.ID == "" && u.Name == "" && u.Password == ""
}

// UserRepository struct representing repository for user data.
type UserRepository struct {
	database *Database
}

// NewUserRepository returns new instance of UserRepository struct
// based on RethinkDB session struct.
func NewUserRepository(database *Database) *UserRepository {
	return &UserRepository{database: database}
}

// SaveUser persists User data into db. Returns an error if
// something bad has happened.
func (repo *UserRepository) SaveUser(user User) error {
	uuid, err := repo.database.createUUID()
	if err != nil {
		return err
	}
	user.ID = uuid

	return repo.collection().Insert(user).Exec(repo.database.Session)
}

// FindUser returns user with given name, 'true' and nil if user exists, empty
// user struct, 'false' and nil if it doesn't exist and error if something
// bad has happened.
func (repo *UserRepository) FindUser(name string) (User, bool, error) {
	user := new(User)

	c, err := repo.collection().Filter(r.Row.Field("name").Eq(name)).Run(repo.database.Session)
	if err != nil {
		return *user, false, err
	}

	if c.IsNil() {
		return *user, false, nil
	}

	if err = c.One(user); err != nil {
		return *user, false, err
	}

	return *user, true, nil
}

func (repo *UserRepository) collection() r.Term {
	return r.DB(repo.database.Config.DBName).Table(repo.database.Config.UsersTableName)
}
