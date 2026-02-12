package validator

import (
	"errors"
	"strconv"
)

const (
	DefaultLimit = 10
	MaxLimit     = 100
)

func ParseLimit(limitStr string) (int, error) {
	if limitStr == "" {
		return DefaultLimit, nil
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return 0, errors.New("limit must be positive number")
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}

	return limit, nil
}
