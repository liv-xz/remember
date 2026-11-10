package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"server/applog"
	"server/config"
	"server/database/connection"
	init_database "server/database/init"
	myredis "server/database/redis"
	"server/flags"
	"server/utils/validate"
	"server/wire"
	"syscall"
)

func main() {

	flags.Parse1()
	fmt.Println("配置文件:", flags.FlagOptions)

	applog.Init()
	applog.Info("日志系统已初始化")

	validate.InitTranslator()

	config.Init(flags.FlagOptions.File)
	applog.Infof("当前环境: %s", config.Global.App.Env)

	uploadDir := config.Global.Upload.Dir
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		applog.Fatalf("创建上传目录失败: %v", err)
	}
	absDir, err := filepath.Abs(uploadDir)
	if err != nil {
		applog.Warnf("解析上传目录绝对路径失败: %v，使用相对路径", err)
		absDir = uploadDir
	}
	applog.Infof("上传目录: %s", absDir)

	if err := connection.Init(); err != nil {
		applog.Fatalf("数据库连接失败: %v", err)
	}
	defer connection.Close()
	applog.Info("数据库连接成功")

	if err := init_database.CreateTable(); err != nil {
		applog.Fatalf("创建数据表失败: %v", err)
	}

	if err := myredis.Init(config.Global.Redis); err != nil {
		applog.Warnf("Redis 连接失败（降级运行）: %v", err)
	} else {
		applog.Info("Redis 连接成功")
		defer myredis.Close()
	}

	engine, err := wire.InitApp()
	if err != nil {
		applog.Fatalf("依赖注入初始化失败: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", config.Global.App.Host, config.Global.App.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	go func() {
		applog.Infof("服务启动: %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			applog.Fatalf("服务异常退出: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	applog.Info("正在关闭服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		applog.Fatalf("服务强制关闭: %v", err)
	}
	applog.Info("服务已关闭")
}
