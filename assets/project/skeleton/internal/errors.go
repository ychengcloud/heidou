package internal

import (
	"{{ . }}/internal/gen/models"
)

const (
	// RetCode ...
	RetCode string = "code"
	// RetMsg ...
	RetMsg string = "msg"
)

const (
	RetCodeOK                  = 2000
	RetCodeUnauthorized        = 4001
	RetCodeNotFound            = 4004
	RetCodeConflict            = 4005
	RetCodeBadParamInput       = 4006
	RetCodeInternalServerError = 5000
)

func GetStatusCode(err error) int {
	if err == nil {
		return RetCodeOK
	}

	switch err {
	case models.ErrInternalServerError:
		return RetCodeInternalServerError
	case models.ErrNotFound:
		return RetCodeNotFound
	case models.ErrConflict:
		return RetCodeConflict
	case models.ErrTokenExpried:
		return RetCodeUnauthorized
	case models.ErrBadPassword:
		return RetCodeUnauthorized
	case models.ErrBadParamInput:
		return RetCodeBadParamInput
	case models.ErrUnauthorized:
		return RetCodeUnauthorized
	default:
		return RetCodeInternalServerError
	}
}
