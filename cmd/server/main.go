// cmd/server/main.go

// 每个可独立执行的Go程序的入口都必须是 main 包
package main

// 导入我们需要的包
import (
	// 标准库
	"log" // 用于打印日志
	"time"

	// 我们自己项目内部的包
	"gotasksys/internal/api/handler"    // 导入我们所有的API处理器 (Handler)
	"gotasksys/internal/api/middleware" // 导入我们所有的中间件 (Middleware)
	"gotasksys/internal/config"         // 导入我们的配置加载和数据库初始化模块

	// 第三方开源库
	"github.com/gin-contrib/cors" // 导入CORS中间件库
	"github.com/gin-gonic/gin"    // 导入Gin框架库
)

// main 函数是整个程序的起点
func main() {
	// --- 步骤1: 加载应用配置 ---
	// 从 config.yaml 文件中读取数据库地址、服务端口等所有配置信息
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		// 如果加载配置失败，程序无法继续，打印致命错误并退出
		log.Fatalf("Failed to load config: %v", err)
	}

	// --- 步骤2: 初始化数据库连接 ---
	// 使用加载到的配置信息，建立与PostgreSQL数据库的连接池
	config.InitDB(cfg)

	// --- 步骤3: 初始化Web框架引擎 ---
	// gin.Default() 会创建一个带有基础中间件（如日志、错误恢复）的Gin引擎
	r := gin.Default()

	// --- 步骤4: 应用全局中间件 ---
	// r.Use(...) 用于向整个应用注册一个或多个中间件，所有请求都会先经过它们
	// cors.Default() 创建一个默认的CORS（跨域资源共享）中间件。
	// 它允许所有源(Origin)的跨域请求，这在前后端分离的本地开发中至关重要。
	// === 用下面这段详细配置，替换掉 r.Use(cors.Default()) ===
	r.Use(cors.New(cors.Config{
		// 允许跨域的源，可以用*通配符，但更安全的方式是明确指定
		AllowOrigins: []string{"http://localhost:5173"},
		// 允许的请求方法
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		// **关键：明确允许 Authorization 这个请求头**
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		// 允许浏览器缓存预检请求结果的时间
		MaxAge: 12 * time.Hour,
	}))
	// ========================================================

	// --- 步骤5: 定义API路由 ---
	// 创建一个API分组，所有接口的URL都会带上 /api/v1 这个前缀，便于版本管理
	v1 := r.Group("/api/v1")
	{
		// 5.1: 定义公开路由 (Public Routes)
		// 这些接口不需要任何认证，任何人都可以访问
		v1.POST("/register", handler.Register) // 管理员创建用户接口 (之前是公开注册)
		v1.POST("/login", handler.Login)       // 用户登录接口

		// 5.2: 定义普通认证路由组 (Authenticated Routes)
		// 创建一个新的子分组，所有在这个组里的路由，都需要先通过它所应用的中间件
		authRequired := v1.Group("/")
		// .Use() 会将中间件应用到这个组内的所有路由上
		// AuthMiddleware 会检查请求头中是否带有合法有效的JWT
		authRequired.Use(middleware.AuthMiddleware())
		{
			// === 新增：为普通用户提供获取任务类型的接口 ===
			authRequired.GET("/task-types", handler.ListTaskTypes) // 获取任务类型列表
			// 所有需要"登录"身份才能访问的接口，都定义在这里
			authRequired.GET("/profile", handler.GetProfile) // 获取个人信息

			// 任务相关的CRUD接口
			authRequired.POST("/tasks", handler.CreateTask)       // 创建任务
			authRequired.GET("/tasks", handler.ListTasks)         // 获取任务列表
			authRequired.GET("/tasks/:id", handler.GetTask)       // 获取单个任务详情
			authRequired.PATCH("/tasks/:id", handler.UpdateTask)  // 更新任务
			authRequired.DELETE("/tasks/:id", handler.DeleteTask) // 删除任务

			// 任务工作流相关的接口
			authRequired.POST("/tasks/:id/approve", handler.ApproveTask)   // 审批任务
			authRequired.POST("/tasks/:id/claim", handler.ClaimTask)       // 领取任务
			authRequired.POST("/tasks/:id/complete", handler.CompleteTask) // 完成任务
			authRequired.POST("/tasks/:id/evaluate", handler.EvaluateTask) // 评价任务

			// 看板和驾驶舱数据接口
			authRequired.GET("/dashboard/summary", handler.GetDashboardSummary) // 获取驾驶舱数据
			authRequired.GET("/personnel/status", handler.GetPersonnelStatus)   // 获取人员看板数据
		}

		// 5.3: 定义管理员路由组 (Admin Routes)
		// 创建一个专门给管理员使用的子分组，路径以 /admin 开头
		adminRoutes := v1.Group("/admin")
		// **注意: 这里应用了两个中间件，它们会按顺序执行**
		// 1. AuthMiddleware 先确保用户已登录
		// 2. AdminAuthMiddleware 再确保该用户角色是 system_admin
		adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
		{
			// 所有只有"系统管理员"才能访问的接口都定义在这里
			adminRoutes.GET("/task-types", handler.ListTaskTypes)   // 获取任务类型列表
			adminRoutes.POST("/task-types", handler.CreateTaskType) // 创建新的任务类型
		}
	}

	// --- 步骤6: 启动HTTP服务 ---
	// 拼接服务地址和端口号
	serverAddr := ":" + cfg.Server.Port
	// 打印一条日志，方便我们知道服务在哪个端口启动
	log.Printf("Server is starting on http://localhost%s", serverAddr)
	// r.Run() 会启动服务并开始监听HTTP请求，它是一个阻塞操作
	if err := r.Run(serverAddr); err != nil {
		// 如果服务启动失败，打印致命错误并退出
		log.Fatalf("Failed to start server: %v", err)
	}
}
