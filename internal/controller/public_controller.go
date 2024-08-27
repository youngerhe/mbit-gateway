package Controller

import (
	"fmt"
	"gateway/internal/constant"
	"gateway/internal/protoc"
	"gateway/pkg/nacos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PublicController struct {
	ServiceName string
	BaseController
}

func (c *PublicController) GetCaptcha(ctx *gin.Context) {

	instance, err := nacos.NacosClient.GetOneHealthyInstance(nacos.Config.App.Services[c.ServiceName])
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", instance.Ip, instance.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
		return
	}
	client := protoc.NewPublicClient(conn)
	rsp, err := client.GetCaptcha(ctx, &protoc.PublicReq{})
	if err != nil {

		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
		return
	}
	c.JsonResp(ctx, int(rsp.Code), rsp.Data)
	return
}

func (c *PublicController) PostCaptcha(ctx *gin.Context) {
	var req *protoc.PostCaptchaReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		c.JsonResp(ctx, constant.ENTITY_ERROR, nil)
		return
	}
	instance, err := nacos.NacosClient.GetOneHealthyInstance(nacos.Config.App.Services[c.ServiceName])
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", instance.Ip, instance.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
		return
	}
	client := protoc.NewPublicClient(conn)
	rsp, err := client.PostCaptcha(ctx, req)
	if err != nil {

		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
		return
	}
	c.JsonResp(ctx, int(rsp.Code), rsp.Data)
	return
}

func (c *PublicController) GetCountry(ctx *gin.Context) {

	instance, err := nacos.NacosClient.GetOneHealthyInstance(nacos.Config.App.Services[c.ServiceName])
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", instance.Ip, instance.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
		return
	}
	client := protoc.NewPublicClient(conn)
	rsp, err := client.GetCountry(ctx, &protoc.PublicReq{})
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
		return
	}
	c.JsonResp(ctx, int(rsp.Code), rsp.Data)
	return
}
