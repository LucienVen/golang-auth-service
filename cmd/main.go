package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/LucienVen/golang-auth-service/internal/app"
	"github.com/LucienVen/golang-auth-service/pkg/log"
)

func main() {
	// 打印当前工作目录
	if wd, err := os.Getwd(); err != nil {
		log.Errorf("获取工作目录失败: %v", err)
	} else {
		log.Infof("当前工作目录: %s", wd)
	}

	// 创建应用实例
	application := app.NewApplication()

	// 启动应用
	if err := application.Start(); err != nil {
		log.Errorf("应用启动失败: %v", err)
		os.Exit(1)
	}

	// 等待中断信号2
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("正在关闭服务...")

	// 创建关闭上下文
	ctx := context.Background()

	// 优雅关闭
	if err := application.Shutdown(ctx); err != nil {
		log.Errorf("服务关闭错误: %v", err)
		os.Exit(1)
	}

	log.Info("服务已关闭")
}
