// Package main 应用入口，启动 HTTP 服务器并监听优雅退出信号
package main

import (
	"backend/bootstrap"
	"backend/pkg/global_vars"
	"backend/pkg/logger"
	"backend/router"
	"context"
	"errors"
	"fmt"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	bootstrap.Bootstrap()

	r := router.SetupRouter()
	addr := global_vars.ConfigYml.GetString("Port")
	server := &http.Server{
		Addr:    addr,
		Handler: r, // 将 gin router 赋值给 Handler
	}
	traceLogger, _ := logger.GetTraceLogger()
	// 2. 在 Goroutine 中启动服务器
	go func() {
		traceLogger.Info(fmt.Sprintf("🚀 服务器正在启动于 %s ...\n", addr))
		// 启动服务
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ 服务器启动失败: %v\n", err)
		}
	}()

	time.Sleep(100 * time.Millisecond) // 等待 server goroutine 就绪

	// ✅ 4. 启动成功回调逻辑
	traceLogger.Info(fmt.Sprintf("✅ 服务器启动成功！ http://localhost%s\n", global_vars.ConfigYml.GetString("Port")))

	// 5. 优雅退出监听
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	traceLogger.Info("\n🛑 正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer global_vars.Db.Close()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}
	traceLogger.Info("服务器已优雅退出")
}
