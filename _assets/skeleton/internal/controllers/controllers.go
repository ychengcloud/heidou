package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"

	"{{ .Extra.pkgpath }}/pkg/transports/http"
	"{{ .Extra.pkgpath }}/internal/services"
)

func CreateInitControllersFn(
	logger *zap.Logger,
) http.InitControllers {
	return func(r *gin.Engine) {
	}
}

var ProviderSet = wire.NewSet(
	CreateInitControllersFn,
)
