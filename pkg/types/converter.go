package types

import (
	"goblog/pkg/logger"
	"strconv"
)

func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func StringToUint64(s string) uint64 {
	i, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		logger.LogError(err)
	}

	return i
}

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.LogError(err)
	}

	return i
}
