// cmd/server/main.go

package main

import (
	"gotasksys/internal/api/handler"
	"gotasksys/internal/api/middleware"
	"gotasksys/internal/config"
	"log"

	// 确保导入handler包
	"github.com/gin-gonic/gin"
)

func main() {
	// ... (加载配置和初始化DB部分不变) ...
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	config.InitDB(cfg)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{

		// --- 公开路由 (不需要登录) ---
		v1.POST("/register", handler.Register)
		v1.POST("/login", handler.Login)
		// --- 受保护路由 (需要登录和JWT) ---
		authRequired := v1.Group("/")
		authRequired.Use(middleware.AuthMiddleware()) // <--- 在这里应用我们的认证中间件
		{
			// 所有需要登录才能访问的接口，都写在这里面
			authRequired.GET("/profile", handler.GetProfile)
			// === 新增的创建任务路由 ===
			authRequired.POST("/tasks", handler.CreateTask)
			// === 新增的获取任务列表路由 ===
			authRequired.GET("/tasks", handler.ListTasks)
			// === 新增的获取单个任务路由 ===
			authRequired.GET("/tasks/:id", handler.GetTask)
			// === 新增的更新任务路由 ===
			// 我们使用PATCH，因为它更符合“部分更新”的语义
			authRequired.PATCH("/tasks/:id", handler.UpdateTask)
			// === 新增的删除任务路由 ===
			authRequired.DELETE("/tasks/:id", handler.DeleteTask)
		}
	}

	// 启动服务
	serverAddr := ":" + cfg.Server.Port
	log.Printf("Server is starting on http://localhost%s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
