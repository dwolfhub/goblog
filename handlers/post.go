package handlers

import (
	"goapi/models"
	"net/http"
)

// GetPostsHandlerFactory creates a func that handles requests to retrieve posts
func GetPostsHandlerFactory(postStore models.PostDataReader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
