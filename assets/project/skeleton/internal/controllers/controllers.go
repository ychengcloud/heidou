package controllers

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"

	"{{ . }}/gen/models"
	"{{ . }}/pkg/auth"
	"{{ . }}/pkg/transports/http"
	jwtmid "{{ . }}/pkg/transports/http/middlewares/jwt"
)

//ParseClaim ...
func parseClaims(c *gin.Context, key string) (userID uint32, userName string, mid uint32, err error) {
	claims, ok := c.Get(key)
	if !ok {
		return 0, "", 0, models.ErrUnauthorized
	}
	claimsMap := claims.(jwt.MapClaims)

	userName = claimsMap["userName"].(string)
	// jwt token中的数值默认序列化成了float64
	userID = uint32(claimsMap["userId"].(float64))
	mid = uint32(claimsMap["merchantsId"].(float64))

	return
}

func CreateInitControllersFn(
	logger *zap.Logger,
	auth *auth.JWTAuth,
	employee *EmployeeController,
	role *RoleController,
	mediaFile *MediaFileController,
) http.InitControllers {
	return func(r *gin.Engine) {
		apiv1 := r.Group("/api/v1")
		// apiv1.Use(jwtmid.New(r, auth, logger).Middleware())
		{
			apiv1.POST("/signin", employee.Signin)
			apiv1.POST("/signout", employee.Signout)
		}
		apiv1.Use(jwtmid.New(r, auth, logger).Middleware())
		{
			apiv1.GET("/current/employee", employee.GetCurrentEmployee)
			apiv1.GET("/current/menus", employee.GetCurrentMenus)
			apiv1.POST("/employee", employee.Create)
			apiv1.POST("/role", role.Create)
			apiv1.PATCH("/role", role.Update)
			apiv1.DELETE("/role/:id", role.Delete)
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
