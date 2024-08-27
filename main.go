package main

import (
	"context"
	"errors"
	"fmt"
	"gateway/internal/router"
	"gateway/pkg/jwt"
	"gateway/pkg/logger"
	"gateway/pkg/nacos"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	nacos.Init()
	nacos.InitConfig()
	logger.Init()
	jwt.Init()

	// 路由初始化
	r := router.Init()

	// 服务启动
	srv := &http.Server{
		Addr:    ":" + viper.GetString("app.port"),
		Handler: r,
	}
	go func() {
		nacos.RegisterInstance()
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {

			nacos.DeregisterInstance()
			panic(err)
		}
	}()

	// -----------------------------优雅退出 -----------------------------
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞
	<-quit

	nacos.DeregisterInstance()
	// 关闭http
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.Error("server shutdown err:", err)
		fmt.Printf("server shutdown: %v ", err)
	}

}
