package services

import "errors"

var (
	errOrderByInvalid = errors.New("order by param is invalid")
)

const (
	ORDER_BY_ASC  string = "ASC"
	ORDER_BY_DESC string = "DESC"
)
