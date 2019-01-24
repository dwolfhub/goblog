package models

// User represents a user object
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
	Active   bool   `json:"-"`
	Created  string `json:"-"`
	Updated  string `json:"-"`
}

// UserDataReader defines methods to retrieve user data
type UserDataReader interface {
	GetUserByUsername(username string) (user User, err error)
	GetUserByEmail(email string) (user User, err error)
}

// GetUserByUsername retrieves a User by username
func (db *DB) GetUserByUsername(username string) (user User, err error) {
	row := db.QueryRow(`
		SELECT id, username, email, password, active, created, updated
		FROM user
		WHERE username = ?
	`, username)

	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Active, &user.Created, &user.Updated)

	return
}

// GetUserByEmail retrieves a User by email
func (db *DB) GetUserByEmail(email string) (user User, err error) {
	row := db.QueryRow(`
		SELECT id, username, email, password, active, created, updated
		FROM user
		WHERE email = ?
	`, email)

	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Active, &user.Created, &user.Updated)

	return
}
