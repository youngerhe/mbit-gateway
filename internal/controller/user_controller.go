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

type UserController struct {
	ServiceName string
	BaseController
}

func (c *UserController) GetEmailCode(ctx *gin.Context) {
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

func (c *UserController) PostEmailCode(ctx *gin.Context) {
	var req *protoc.PostEmailCodeReq
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
	client := protoc.NewUcClient(conn)
	rsp, err := client.PostEmailCode(ctx, req)
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
		return
	}
	c.JsonResp(ctx, int(rsp.Code), rsp.Data)
	return
}

func (c *UserController) Register(ctx *gin.Context) {
	var req *protoc.RegisterReq
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
	client := protoc.NewUcClient(conn)
	rsp, err := client.Register(ctx, req)
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
		return
	}
	c.JsonResp(ctx, int(rsp.Code), rsp.Data)
	return
}

func (c *UserController) Login(ctx *gin.Context) {
	var req *protoc.LoginReq
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
	client := protoc.NewUcClient(conn)
	rsp, err := client.Login(ctx, req)
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
		return
	}
	c.JsonResp(ctx, int(rsp.Code), rsp.Data)
	return
}

func (c *UserController) GetUserInfo(ctx *gin.Context) {
	instance, err := nacos.NacosClient.GetOneHealthyInstance(nacos.Config.App.Services[c.ServiceName])
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", instance.Ip, instance.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
		return
	}
	client := protoc.NewUcClient(conn)

	rsp, err := client.GetUserInfo(ctx, &protoc.GetUserInfoReq{Uid: ctx.GetInt64("uid")})
	if err != nil {
		c.JsonResp(ctx, constant.SYSTEM_ERROR, nil)
		return
	}
	c.JsonResp(ctx, int(rsp.Code), rsp.Data)
	return
}
