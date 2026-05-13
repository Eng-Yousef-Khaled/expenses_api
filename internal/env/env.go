package env

import "os"

func GetString(key string, fullback string) string {
	if res := os.Getenv(key); res != "" {
		return res
	}
	return fullback
}
