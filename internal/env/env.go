package env

import (
	"os"
	"strconv"
)

// GetString retrieves the string value for an environment
// variable corresponding to the given key,
// or the fallback if no matching key is present.
func GetString(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

// GetInt retrieves the integer value for an environment
// variable corresponding to the given key,
// or the fallback if no matching key is present
// or in the case of error converting to integer.
func GetInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return result
}
