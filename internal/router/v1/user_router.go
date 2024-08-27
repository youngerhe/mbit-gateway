package v1

import (
	Controller "gateway/internal/controller"
	"gateway/internal/router/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(e *gin.Engine) {
	r := e.Group("/uc/v1/")
	{
		uc := Controller.UserController{ServiceName: "uc"}
		r.GET("/email/code", uc.GetEmailCode)
		r.POST("/email/code", uc.PostEmailCode)
		r.POST("/register", uc.Register)
		r.POST("/login", uc.Login)
		// 需要鉴权的接口
		ur := r.Use(middleware.JwtMiddleware())
		{
			ur.GET("/user", uc.GetUserInfo)
		}
	}
}
