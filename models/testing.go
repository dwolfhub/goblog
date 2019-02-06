package models

// MockUserDataStore can be used to mock IUserDataStore
type MockUserDataStore struct {
	User User
	Err  error
}

// GetUserByUsername retrieves user and error
func (m *MockUserDataStore) GetUserByUsername(username string) (User, error) {
	return m.User, m.Err
}

// GetUserByEmail retrieves user and error
func (m *MockUserDataStore) GetUserByEmail(email string) (User, error) {
	return m.User, m.Err
}

// GetMockUserDataStore is a helper method to retrieve a MockUserDataStore
func GetMockUserDataStore(err error, user *User) (ud *MockUserDataStore) {
	ud = &MockUserDataStore{}
	ud.Err = err
	if user != nil {
		ud.User = *user
	}

	return
}
