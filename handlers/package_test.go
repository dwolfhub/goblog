package handlers_test

import (
	"bytes"
	"goapi/models"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

func testHttpRequest(body string, handler func(w http.ResponseWriter, r *http.Request)) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	method := "GET"
	if len(body) > 0 {
		method = "POST"
	}

	req, _ := http.NewRequest(method, "/", bytes.NewBufferString(body))

	r := mux.NewRouter()
	r.HandleFunc("/", handler)

	r.ServeHTTP(rr, req)

	return rr
}

type mockUserDataStore struct {
	user models.User
	err  error
}

func (m *mockUserDataStore) GetUserByUsername(username string) (models.User, error) {
	return m.user, m.err
}
func (m *mockUserDataStore) GetUserByEmail(email string) (models.User, error) {
	return m.user, m.err
}

func getMockUserDataStore(err error, user *models.User) (ud *mockUserDataStore) {
	ud = &mockUserDataStore{}
	ud.err = err
	ud.user = *user

	return
}

type mockEmailSender struct {
	f   func()
	err error
}

func (m *mockEmailSender) EmailSend(toAddresses []string, subject string, body string) error {
	m.f()
	return m.err
}

func getMockEmailSender(err error, f func()) (es *mockEmailSender) {
	es = &mockEmailSender{}
	es.f = f

	return
}
