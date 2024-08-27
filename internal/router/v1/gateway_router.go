package v1

import (
	Controller "gateway/internal/controller"
	"github.com/gin-gonic/gin"
)

// GatewayRouter 自身提供的接口
func GatewayRouter(e *gin.Engine) {
	r := e.Group("/v1/")
	{
		gateway := Controller.GatewayController{}
		r.GET("/token/refresh", gateway.RefreshToken)
	}
}
