package misc

import (
	"os"
	"strconv"
)

func GetEnv(key, _default string) string {
	data := os.Getenv(key)
	if len(data) > 0 {
		return data
	}
	return _default
}

func GetIntEnv(key string, _default int) int {
	data := os.Getenv(key)
	if len(data) > 0 {
		intData, err := strconv.Atoi(data)
		if err != nil {
			return _default
		}
		return intData
	}
	return _default
}

func GetEnvOrPanic(key string) string {
	data := os.Getenv(key)
	if len(data) > 0 {
		return data
	}
	panic("Не указана переменная окружения: " + key)
}
