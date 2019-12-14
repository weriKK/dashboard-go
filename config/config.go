package config

import "os"

import "strconv"

// GetString returns an environment variable as string
func GetString(key string) string {

	return os.Getenv(key)
}

// GetInt returns an environment variable as int. If the value can not be parsed as an int, returns 0.
func GetInt(key string) int {

	strVal := os.Getenv(key)

	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		return 0
	}

	return intVal
}
