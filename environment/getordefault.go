package environment

import (
	"os"
)

func GetOrDefault(key, defaultvalue string) string {
	if val := os.Getenv(key); val == "" {
		return defaultvalue
	} else {
		return val
	}
}
