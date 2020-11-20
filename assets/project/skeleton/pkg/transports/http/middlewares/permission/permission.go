package jwt

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"

	"{{ . }}/pkg/auth"
)

type Callback interface {
	Unauthorized(c *gin.Context, err error)
}

// GinPermission  struct
type GinPermission struct {
	auth     *auth.JWTAuth
	logger   *zap.Logger
	callback Callback
}

// Init 初始化数据库
func New(a *auth.JWTAuth, logger *zap.Logger, cb Callback) *GinPermission {
	ginPermission := &GinPermission{
		auth:     a,
		logger:   logger,
		callback: cb,
	}

	return ginPermission
}

// GinPermission is permission middleware
func (j *GinPermission) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := j.auth.Authorization(c)
		if err != nil {
			j.callback.Unauthorized(c, err)

			c.Abort()
			return
		}

		if !ok {
			j.callback.Unauthorized(c, fmt.Errorf("Unthorization"))

			c.Abort()
			return
		}

		c.Next()
	}
}

var ProviderSet = wire.NewSet(New)
