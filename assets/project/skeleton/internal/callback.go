package internal

import (
	"github.com/gin-gonic/gin"

	"{{ . }}/pkg/transports/http/middlewares/jwt"
	permission "{{ . }}/pkg/transports/http/middlewares/permission"
)

type JWTCallback struct {
}

func NewJWTCallback() jwt.Callback {
	return &JWTCallback{}
}
func (cb *JWTCallback) Unauthentication(c *gin.Context, err error) {
	c.AbortWithStatusJSON(200, gin.H{RetCode: GetStatusCode(err), RetMsg: err.Error()})
}

type AuthCallback struct {
}

func NewAuthCallback() permission.Callback {
	return &AuthCallback{}
}

func (cb *AuthCallback) Unauthorized(c *gin.Context, err error) {
	c.AbortWithStatusJSON(200, gin.H{RetCode: GetStatusCode(err), RetMsg: err.Error()})
}
