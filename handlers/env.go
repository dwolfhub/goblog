package handlers

import (
	"goapi/models"
)

// Env holds instances of services and defines route handlers
type Env struct {
	DB models.IDataStore
}
