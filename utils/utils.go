package utils

import (
	"os"
	"strconv"
)

func GetEnv(key string, defaultValue string) string {
	val := os.Getenv(key)

	if val == "" {
		return defaultValue
	}
	return val
}

func Atoi(s string, defaultValue int) int {
	n, err := strconv.Atoi(s)

	if err != nil {
		return defaultValue
	}
	return n
}

func Atoi64(s string, defaultValue int64) int64 {
	n, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return defaultValue
	}
	return n
}
func Atof(s string, defaultValue float64) float64 {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return n
}
