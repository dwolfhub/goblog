package handlers

import (
	"encoding/json"
	"goapi/models"

	"github.com/gorilla/context"

	"net/http"
)

// GetMeHandlerFactory creates a func that handles /user/me requests
func GetMeHandlerFactory(userStore models.UserDataReader) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if username, ok := context.Get(req, "username").(string); ok {
			if user, err := userStore.GetUserByUsername(username); err == nil {
				json.NewEncoder(w).Encode(user)
				return
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
