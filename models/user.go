package models

// User represents a user object
type User struct {
	ID       int
	Username string
	Password string
	Email    string
	Active   bool
	Created  string
	Updated  string
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
