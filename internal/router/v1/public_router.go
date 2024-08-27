package v1

import (
	Controller "gateway/internal/controller"
	"github.com/gin-gonic/gin"
)

// 公共服务
func PublicRouter(e *gin.Engine) {
	r := e.Group("/public/v1/")
	{
		public := Controller.PublicController{ServiceName: "uc"}
		r.GET("/captcha", public.GetCaptcha)
		r.POST("/captcha", public.PostCaptcha)
		r.GET("/country", public.GetCountry)
	}
}
