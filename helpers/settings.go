package helpers

import (
	"log"
	"os"
)

// GetEnvVar retrieves an environment variable or dies trying
func GetEnvVar(key string) string {
	val, isSet := os.LookupEnv(key)

	if !isSet {
		log.Fatalf("Environment variable %s must be set.", key)
	}

	return val
}
