package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/listener"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)
import _ "dros-go-plugin/docs"
import _ "dros-go-plugin/plugins/api_report"
import _ "dros-go-plugin/plugins/service_register"
import ginSwagger "github.com/swaggo/gin-swagger"
import swaggerFiles "github.com/swaggo/files"

var engine *gin.Engine = nil

// @title           dros-go-plugin
// @version         1.0
// @description     插件管理api
// @tag.name 插件管理
// @tag.description 插件管理测试api
// @tag.name 任务管理
// @tag.description 任务管理测试api
// @host      localhost:8080
// @BasePath  /api/dros-go-plugin

// @securityDefinitions.basic  BasicAuth
func main() {
	defer func() {
		fmt.Println("异常退出")
		if err := recover(); err != nil {
			fmt.Printf("异常原因:%v\n", err)
			return
		}
	}()

	engine = gin.New()
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	serverStart()
}

func serverStart() {
	listener.PublishEvent(listener.ServerRunStartEvent{})
	go func() {
		err := engine.Run(":8080")
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("服务启动异常")
		} else {
			listener.PublishEvent(listener.ServerStopEvent{})
		}
		fmt.Println("服务启动成功 ")
	}()
	listener.PublishEvent(listener.ServerRunFinishEvent{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
