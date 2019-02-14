package errors

import (
	"encoding/json"
)

// Error ...
type Error struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// New 常规业务逻辑错误
func New(code int, detail string) error {
	return &Error{
		Code:   code,
		Detail: detail,
	}
}

// Unauthorized 未经授权错误
func Unauthorized(detail string) error {
	return &Error{
		Code:   401,
		Detail: detail,
	}
}

// Internal 内部错误
func Internal(err error) error {
	return &Error{
		Code:   500,
		Detail: err.Error(),
	}
}
