package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func GetStringOrDefault(key string, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	conv, err := strconv.Atoi(val)
	if err != nil {
		//TODO: Perhaps log something
		return fallback
	}
	return conv
}

func GetString(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", errors.New(fmt.Sprintf("Failed to fetch env var %s", key))
	}
	return val, nil
}

func IsProd() bool {
	val, err := GetString("APP_STAGE")
	if err != nil || val != "PROD" {
		return false
	}
	return true
}
