package db

// User is a struct containing user data.
type User struct {
	ID       string `json:"id" gorethink:"id,omitempty"`
	Login    string `json:"login" gorethink:"name,omitempty"`
	Password string `json:"password" gorethink:"password,omitempty"`
}

// Empty returns 'true' it the User struct is empty, false otherwise.
func (u *User) Empty() bool {
	return u.ID == "" && u.Login == "" && u.Password == ""
}

func (u *User) Name() string {
	return u.Login
}

type Table interface {
	UUID() (string, error)
	Insert(interface{}) error
	Find(property string, value interface{}, result interface{}) error
}

// UserRepository struct representing repository for user data.
type UserRepository struct {
	table Table
}

// NewUserRepository returns new instance of UserRepository struct
// based on RethinkDB session struct.
func NewUserRepository(table Table) *UserRepository {
	return &UserRepository{table: table}
}

// SaveUser persists User data into db. Returns an error if
// something bad has happened.
func (r *UserRepository) SaveUser(user User) error {
	uuid, err := r.table.UUID()
	if err != nil {
		return err
	}
	user.ID = uuid
	return r.table.Insert(user)
	// repo.database.Session.GetTable(repo.database.Config.UsersTableName).Insert(user)
	//return repo.collection().Insert(user).Exec(repo.database.Session)
}

// FindUser returns user with given name, 'true' and nil if user exists, empty
// user struct, 'false' and nil if it doesn't exist and error if something
// bad has happened.
func (r *UserRepository) FindUser(name string) (User, error) {
	user := new(User)
	if err := r.table.Find("name", name, user); err != nil {
		return *user, err
	}
	return *user, nil
}
