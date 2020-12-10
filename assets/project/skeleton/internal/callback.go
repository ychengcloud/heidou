package internal

import (
	"github.com/gin-gonic/gin"



	permission "{{ . }}/pkg/transports/http/middlewares/permission"
	"{{ . }}/pkg/transports/http/middlewares/jwt"
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
