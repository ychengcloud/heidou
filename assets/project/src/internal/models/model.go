package models

import (
	"errors"
)

var (
	// ErrNameServed ...
	ErrNameServed = errors.New("NameServed")
	ErrSignFail   = errors.New("SignFail")
)

// LoginTokenInfo 登录令牌信息
type LoginTokenInfo struct {
	AccessToken string `json:"accessToken" swaggo:"true,访问令牌"`
	TokenType   string `json:"tokenType" swaggo:"true,令牌类型"`
	ExpiresAt   int64  `json:"expiresAt" swaggo:"true,令牌到期时间"`
}
