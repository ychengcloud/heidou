package controllers

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"

	"{{ . }}/internal/gen/models"
	"{{ . }}/pkg/auth"
	"{{ . }}/pkg/transports/http"
	jwtmid "{{ . }}/pkg/transports/http/middlewares/jwt"
)

//ParseClaim ...
func parseClaims(c *gin.Context, key string) (userID uint64, userName string, mid uint64, err error) {
	claims, ok := c.Get(key)
	if !ok {
		return 0, "", 0, models.ErrUnauthorized
	}
	claimsMap := claims.(jwt.MapClaims)

	userName = claimsMap[auth.ClaimsUsernameKey].(string)
	// jwt token中的数值默认序列化成了float64
	userID = uint64(claimsMap[auth.ClaimsUserIDKey].(float64))
	mid = uint64(claimsMap[auth.ClaimsMerchantsIDKey].(float64))

	return
}

func CreateInitControllersFn(
	logger *zap.Logger,
	ginJWT *jwtmid.GinJWT,
	auth *auth.JWTAuth,
	employee *EmployeeController,
	role *RoleController,
	mediaFile *MediaFileController,
) http.InitControllers {
	return func(r *gin.Engine) {
		apiv1 := r.Group("/api/v1")
		{
			apiv1.POST("/signin", employee.Signin)
		}
		apiv1.Use(ginJWT.Middleware())
		{
			apiv1.POST("/signout", employee.Signout)
			apiv1.GET("/current/employees", employee.GetCurrentEmployee)
			apiv1.GET("/current/menus", employee.GetCurrentMenus)
			apiv1.POST("/employees", employee.Create)
			apiv1.POST("/roles", role.Create)
			apiv1.PATCH("/roles", role.Update)
			apiv1.DELETE("/roles/:id", role.Delete)
			apiv1.POST("/upload", mediaFile.Upload)
			apiv1.GET("/upload/sign", mediaFile.GetSign)
		}
	}
}

var ProviderSet = wire.NewSet(
	CreateInitControllersFn,
	NewEmployeeController,
	NewRoleController,
	NewMediaFileController,
)
