package tools

import "errors"

var (
	errIDRequired    = errors.New("id is required")
	errQueryRequired = errors.New("query is required")
)
