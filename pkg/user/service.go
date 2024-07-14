package user

const nameProp = "name"

type Database interface {
	UUID() (string, error)
	Insert(interface{}) error
	Find(property string, value interface{}, result interface{}) error
}

// Service struct representing repository for user data.
type Service struct {
	db Database
}

// NewUserService returns new instance of Service.
func NewUserService(db Database) *Service {
	return &Service{db: db}
}

// SaveUser persists User data into db. Returns an error if
// something bad has happened.
func (s *Service) SaveUser(user User) error {
	uuid, err := s.db.UUID()
	if err != nil {
		return err
	}

	user.ID = uuid

	return s.db.Insert(user)
}

// FindUser returns user with given name
func (s *Service) FindUser(name string) (*User, error) {
	var user User
	if err := s.db.Find(nameProp, name, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
