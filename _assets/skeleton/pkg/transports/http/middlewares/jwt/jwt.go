package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"

	"{{ .Extra.pkgpath }}/pkg/auth"
)

type Callback interface {
	Unauthentication(c *gin.Context, err error)
}

// GinJWT  struct
type GinJWT struct {
	auth     *auth.JWTAuth
	logger   *zap.Logger
	callback Callback
}

// Init 初始化数据库
func New(a *auth.JWTAuth, logger *zap.Logger, cb Callback) *GinJWT {
	ginJWT := &GinJWT{
		auth:     a,
		logger:   logger,
		callback: cb,
	}

	return ginJWT
}

// JWT is jwt middleware
func (j *GinJWT) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := j.auth.Authentication(c)
		if err != nil {
			j.callback.Unauthentication(c, err)

			c.Abort()
			return
		}

		c.Set(j.auth.ClaimsKey, claims)

		c.Next()
	}
}

var ProviderSet = wire.NewSet(New)
