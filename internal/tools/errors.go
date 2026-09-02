package tools

import "errors"

var (
	errIDRequired           = errors.New("id is required")
	errQueryRequired        = errors.New("query is required")
	errNameRequired         = errors.New("name is required")
	errTitleRequired        = errors.New("title is required")
	errRefCodeRequired      = errors.New("ref_code is required")
	errUpdateFieldsRequired = errors.New("at least one field to update is required")
)
