// cmd/server/main.go

package main

import (
	"gotasksys/internal/config" // 导入我们自己的包
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	config.InitDB(cfg)

	// 初始化 Gin 引擎
	r := gin.Default()

	// 设置一个基础的 API 分组
	v1 := r.Group("/api/v1")
	{
		// 健康检查路由
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// 启动服务
	serverAddr := ":" + cfg.Server.Port
	log.Printf("Server is starting on http://localhost%s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
