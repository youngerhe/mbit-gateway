package router

import (
	"gateway/internal/router/v1"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()
	v1.UserRouter(r)
	v1.PublicRouter(r)
	v1.GatewayRouter(r)
	return r
}
