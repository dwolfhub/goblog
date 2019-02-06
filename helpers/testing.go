package helpers

// MockEmailSender mocks IEmailSender
type MockEmailSender struct {
	F   func()
	Err error
}

// EmailSend runs method and then returns error
func (m *MockEmailSender) EmailSend(toAddresses []string, subject string, body string) error {
	m.F()
	return m.Err
}

// GetMockEmailSender is a helper method to retrieve a MockEmailSender
func GetMockEmailSender(err error, f func()) (es *MockEmailSender) {
	es = &MockEmailSender{}
	es.F = f

	return
}
