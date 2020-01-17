package user

// User is a struct containing user data.
type User struct {
	ID       string `json:"id" gorethink:"id,omitempty"`
	Login    string `json:"login" gorethink:"name,omitempty"`
	Password string `json:"password" gorethink:"password,omitempty"`
}

// Empty returns 'true' it the User struct is empty, false otherwise.
func (u *User) Empty() bool {
	return u == nil || (u.ID == "" && u.Login == "" && u.Password == "")
}

func (u *User) Name() string {
	return u.Login
}
