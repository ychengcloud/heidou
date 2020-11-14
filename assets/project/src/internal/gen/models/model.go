package models

import (
	"errors"
	"fmt"
)

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("Internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("Your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("Your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("Given Param is not valid")
	// ErrBadPassword will throw if the given password is wrong
	ErrBadPassword = errors.New("Given password is wrong")
	// ErrTokenExpried Token expried
	ErrTokenExpried = errors.New("Token expried")
	// ErrUnauthorized Unauthorized
	ErrUnauthorized = errors.New("Unauthorized")
)

//PaginationQuery gin handler query binding struct
type PaginationQuery struct {
	Where  string `form:"where"`
	Fields string `form:"fields"`
	Order  string `form:"order"`
	Offset uint   `form:"offset"`
	Limit  uint   `form:"limit"`
}

//String to string
func (pq *PaginationQuery) String() string {
	return fmt.Sprintf("w=%v_f=%s_o=%s_of=%d_l=%d", pq.Where, pq.Fields, pq.Order, pq.Offset, pq.Limit)
}
