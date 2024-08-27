package Controller

import (
	"gateway/internal/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct{}

func (c BaseController) JsonResp(ctx *gin.Context, errCode int, data interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	resp := gin.H{
		"code": errCode,
		"msg":  constant.CodeMap[errCode],
		"data": data,
	}
	ctx.JSON(http.StatusOK, resp)
	return
}
