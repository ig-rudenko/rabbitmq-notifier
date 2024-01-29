package misc

import (
	"os"
)

func GetEnv(key, _default string) string {
	data := os.Getenv(key)
	if len(data) > 0 {
		return data
	}
	return _default
}
