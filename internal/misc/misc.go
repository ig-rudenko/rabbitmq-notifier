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
			return intData
		}
	}
	return _default
}

func GetUIntEnv(key string, _default uint64) uint64 {
	data := os.Getenv(key)
	if len(data) > 0 {
		intData, err := strconv.ParseUint(data, 10, 32)
		if err != nil {
			return intData
		}
	}
	return _default
}
