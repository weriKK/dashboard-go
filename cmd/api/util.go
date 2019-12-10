package main

import (
	"errors"
	"net/url"
	"strconv"
)

func GetLimitQueryParam(val url.Values) (int, error) {

	if val.Get("limit") == "" || len(val.Get("limit")) < 1 {
		return 0, errors.New("no 'limit' query param")
	}

	limit, err := strconv.Atoi(val.Get("limit"))
	if err != nil {
		return 0, err
	}

	return limit, nil
}
