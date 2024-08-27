package Controller

import (
	"gateway/internal/constant"
	"gateway/internal/types"
	"gateway/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type GatewayController struct {
	BaseController
}

func (c *GatewayController) RefreshToken(ctx *gin.Context) {

	// 获取refreshToken
	refreshToken := ctx.GetHeader("refresh-token")
	newAccessToken, newRefreshToken, err := jwt.RefreshToken(refreshToken)
	if err != nil {
		c.JsonResp(ctx, constant.REFRESH_TOKEN_FAILED, nil)
		return
	}
	result := types.RefreshTokenResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}
	c.JsonResp(ctx, constant.SUCCESS, result)
}
