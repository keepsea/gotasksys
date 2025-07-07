// cmd/server/main.go
// 每个可独立执行的Go程序的入口都必须是 main 包
package main

// 导入我们需要的包
import (
	// 标准库
	"log" // 用于打印日志
	"time"

	// 项目内部的包
	"gotasksys/internal/api/handler"    // 导入所有的API处理器 (Handler)
	"gotasksys/internal/api/middleware" // 导入所有的中间件 (Middleware)
	"gotasksys/internal/config"         // 导入配置加载和数据库初始化模块

	"github.com/gin-contrib/cors" // 导入CORS中间件库
	"github.com/gin-gonic/gin"    // 导入Gin框架库
)

// main 函数是整个程序的起点

func main() {
	// --- 加载配置、初始化数据库、应用CORS中间件  ---
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	config.InitDB(cfg)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST"}, // 只允许GET和POST
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	// --- 核心：定义所有API路由 ---
	v1 := r.Group("/api/v1")
	{
		// 1. 登录路由(公开路由)
		v1.POST("/login", handler.Login)

		// 2. 用户功能路由组(需要用户认证)
		authRequired := v1.Group("/")
		authRequired.Use(middleware.AuthMiddleware())
		{
			// 用户个人信息相关路由
			authRequired.GET("/profile", handler.GetProfile)                      // 获取当前登录用户的个人信息
			authRequired.POST("/profile/update-details", handler.UpdateMyProfile) // 更新个人信息
			authRequired.POST("/profile/change-password", handler.ChangeMyPassword)
			authRequired.GET("/task-types", handler.ListTaskTypes)
			authRequired.GET("/dashboard/summary", handler.GetDashboardSummary)
			authRequired.GET("/personnel/status", handler.GetPersonnelStatus)

			// === 新增：用户个人请假管理路由 ===
			leaveRoutes := authRequired.Group("/profile/leaves")
			{
				leaveRoutes.GET("", handler.ListMyLeavesHandler)            // 查看我的请假
				leaveRoutes.POST("", handler.CreateLeaveHandler)            // 新增请假
				leaveRoutes.POST("/:id/delete", handler.DeleteLeaveHandler) // 删除请假
			}

			// 任务管理路由(需要用户认证)
			authRequired.POST("/tasks", handler.CreateTask)
			authRequired.GET("/tasks", handler.ListTasks)
			authRequired.GET("/tasks/:id", handler.GetTask)
			authRequired.POST("/tasks/:id/update", handler.UpdateTask)
			authRequired.POST("/tasks/:id/delete", handler.DeleteTask)

			// 任务工作流
			authRequired.POST("/tasks/:id/approve", handler.ApproveTask)
			authRequired.POST("/tasks/:id/reject", handler.RejectTask)
			authRequired.POST("/tasks/:id/resubmit", handler.ResubmitTask)
			authRequired.POST("/tasks/:id/claim", handler.ClaimTask)
			authRequired.POST("/tasks/:id/complete", handler.CompleteTask)
			authRequired.POST("/tasks/:id/evaluate", handler.EvaluateTask)

			// 任务转交
			authRequired.POST("/tasks/:id/transfer", handler.InitiateTransfer)
			authRequired.POST("/transfers/:transfer_id/accept", handler.AcceptTransfer)
			authRequired.POST("/transfers/:transfer_id/reject", handler.RejectTransfer)
			authRequired.POST("/transfers/:transfer_id/cancel", handler.CancelTransfer)

			// 子任务管理路由
			authRequired.POST("/tasks/:id/subtasks", handler.CreateSubtask)

			// 管理员指派任务
			authRequired.POST("/tasks/:id/assign", handler.AssignTask)
			// 用户获取可用头像列表的路由
			authRequired.GET("/system-avatars", handler.ListAvailableAvatars)

			// 计划任务管理路由 (仅Manager可访问)
			periodicRoutes := authRequired.Group("/periodic-tasks")
			// 这里加一个Manager的中间件
			periodicRoutes.Use(middleware.ManagerAuthMiddleware())
			{
				periodicRoutes.GET("", handler.ListPeriodicTasks)
				periodicRoutes.POST("", handler.CreatePeriodicTask)
				periodicRoutes.POST("/:id/update", handler.UpdatePeriodicTask)
				periodicRoutes.POST("/:id/delete", handler.DeletePeriodicTask)
				periodicRoutes.POST("/:id/toggle", handler.TogglePeriodicTask)
			}
		}

		// 3. 管理员路由组
		adminRoutes := v1.Group("/admin")
		adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
		{
			// 用户管理
			adminRoutes.GET("/users", handler.ListUsers)
			adminRoutes.POST("/users", handler.CreateUser)
			adminRoutes.POST("/users/:id/update", handler.UpdateUser)
			adminRoutes.POST("/users/:id/reset-password", handler.ResetPassword)
			adminRoutes.POST("/users/:id/delete", handler.DeleteUser)

			// 任务类型管理
			adminRoutes.GET("/task-types", handler.ListTaskTypes)
			adminRoutes.POST("/task-types", handler.CreateTaskType)
			adminRoutes.POST("/task-types/:id/update", handler.UpdateTaskType)
			adminRoutes.POST("/task-types/:id/delete", handler.DeleteTaskType)
			// 管理员头像库管理路由组
			avatarRoutes := adminRoutes.Group("/system-avatars")
			{
				avatarRoutes.GET("", handler.ListAllSystemAvatars)           // 获取所有头像（包括禁用的）
				avatarRoutes.POST("", handler.CreateSystemAvatar)            // 新增一个头像
				avatarRoutes.POST("/:id/update", handler.UpdateSystemAvatar) // 修改一个头像
				avatarRoutes.POST("/:id/delete", handler.DeleteSystemAvatar) // 删除一个头像
			}
		}
	}

	// --- 启动服务 ---
	serverAddr := ":" + cfg.Server.Port
	log.Printf("Server is starting on http://localhost%s", serverAddr) // 打印服务器启动地址
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
