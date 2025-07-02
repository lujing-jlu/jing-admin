package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User模型现在在models.go中定义

// 数据库实例
var db *gorm.DB

// 初始化数据库
func initDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open("jing_admin.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 自动迁移数据库
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化权限系统
	err = initPermissionSystem()
	if err != nil {
		log.Fatal("Failed to initialize permission system:", err)
	}

	// 初始化日志系统
	err = initLogSystem()
	if err != nil {
		log.Fatal("Failed to initialize log system:", err)
	}

	// 初始化系统配置
	err = initSystemConfig()
	if err != nil {
		log.Fatal("Failed to initialize system config:", err)
	}

	// 初始化文件上传系统
	err = initUploadSystem()
	if err != nil {
		log.Fatal("Failed to initialize upload system:", err)
	}

	// 创建默认管理员账户
	var adminUser User
	result := db.Where("username = ?", "admin").First(&adminUser)
	if result.Error != nil {
		// 加密密码
		hashedPassword, err := hashPassword("admin123")
		if err != nil {
			log.Fatal("Failed to hash admin password:", err)
		}
		
		// 创建默认管理员
		admin := User{
			Username: "admin",
			Email:    "admin@jingadmin.com",
			Password: hashedPassword,
			Role:     "admin",
			Status:   true,
		}
		db.Create(&admin)
		log.Println("Default admin user created: username=admin, password=admin123")
	}

	log.Println("Database initialized successfully")
}

// CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// 日志中间件
func loggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// API响应格式
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 成功响应
func successResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// 错误响应
func errorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, ApiResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func main() {
	// 初始化数据库
	initDatabase()

	// 创建默认的gin路由器
	r := gin.Default()

	// 添加中间件
	r.Use(corsMiddleware())
	r.Use(loggerMiddleware())
	r.Use(gin.Recovery()) // 恢复中间件

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		successResponse(c, gin.H{
			"status":    "ok",
			"message":   "Jing Admin Backend is running",
			"database":  "connected",
			"timestamp": time.Now().Unix(),
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 基础测试接口
		api.GET("/test", func(c *gin.Context) {
			successResponse(c, gin.H{
				"message": "API is working",
				"version": "v1.0.0",
			})
		})

		// 认证相关接口（无需认证）
		auth := api.Group("/auth")
		{
			auth.POST("/login", login)
			auth.POST("/register", register)
			auth.POST("/logout", logout)
		}

		// 公开配置接口（无需认证）
		api.GET("/config/public", getPublicSystemConfigs)

		// 需要认证的接口
		protected := api.Group("/")
		protected.Use(authMiddleware())
		protected.Use(logMiddleware()) // 添加日志中间件
		{
			// 当前用户信息
			protected.GET("/me", getCurrentUser)
			protected.PUT("/me", updateProfile)
			protected.POST("/change-password", changePassword)
			protected.GET("/my-permissions", getUserPermissionsAPI)

			// 用户相关接口（需要管理员权限）
			users := protected.Group("/users")
			users.Use(adminMiddleware())
			{
				users.GET("", getUserList)
				users.GET("/:id", getUserById)
				users.POST("", createUser)
				users.PUT("/:id", updateUser)
				users.DELETE("/:id", deleteUser)
				users.POST("/:id/roles", assignUserRoles)
			}

			// 角色管理接口（需要管理员权限）
			roles := protected.Group("/roles")
			roles.Use(adminMiddleware())
			{
				roles.GET("", getRoleList)
				roles.GET("/:id", getRoleById)
				roles.POST("", createRole)
				roles.PUT("/:id", updateRole)
				roles.DELETE("/:id", deleteRole)
			}

			// 权限管理接口（需要管理员权限）
			permissions := protected.Group("/permissions")
			permissions.Use(adminMiddleware())
			{
				permissions.GET("", getPermissionList)
				permissions.POST("/assign", assignRolePermissions)
			}

			// 操作日志接口（需要管理员权限）
			logs := protected.Group("/logs")
			logs.Use(adminMiddleware())
			{
				logs.GET("", getOperationLogs)
				logs.GET("/:id", getOperationLogById)
				logs.DELETE("/:id", deleteOperationLog)
				logs.POST("/batch-delete", batchDeleteOperationLogs)
				logs.DELETE("/clear-old", clearOldOperationLogs)
				logs.GET("/stats", getOperationLogStats)
			}

			// 系统信息接口
			system := protected.Group("/system")
			{
				system.GET("/info", getSystemInfo)
			}

			// 系统配置接口（需要管理员权限）
			config := protected.Group("/config")
			config.Use(adminMiddleware())
			{
				config.GET("", getSystemConfigs)
				config.GET("/:key", getSystemConfigByKey)
				config.PUT("/:key", updateSystemConfig)
				config.POST("/batch", batchUpdateSystemConfigs)
				config.POST("", createSystemConfig)
				config.DELETE("/:key", deleteSystemConfig)
				config.POST("/:key/reset", resetSystemConfigToDefault)
			}

			// 文件上传接口
			files := protected.Group("/files")
			{
				files.POST("/upload", uploadFile)
				files.GET("", getFileList)
				files.GET("/stats", getFileStats)
				files.DELETE("/:id", deleteFile)
			}

			// 文件访问接口（公开）
			api.GET("/uploads/:filename", serveFile)

			// 数据导出接口（仅管理员）
			protected.GET("/export/users", adminMiddleware(), exportUsersCSV)
			protected.GET("/export/roles", adminMiddleware(), exportRolesCSV)
			protected.GET("/export/permissions", adminMiddleware(), exportPermissionsCSV)

			// 数据导入接口（仅管理员）
			protected.POST("/import/users", adminMiddleware(), importUsersCSV)
			protected.POST("/import/roles", adminMiddleware(), importRolesCSV)
			protected.POST("/import/permissions", adminMiddleware(), importPermissionsCSV)
		}

		// 系统监控API
		system := api.Group("/system")
		{
			system.GET("/monitor", getSystemMonitor)
		}
	}

	log.Println("启动服务器在端口 :8081")
	log.Println("数据库文件: jing_admin.db")
	log.Println("健康检查: http://localhost:8081/health")
	log.Println("API测试: http://localhost:8081/api/test")
	
	r.Run(":8081")
}

// 用户管理API处理函数

// 获取用户列表
func getUserList(c *gin.Context) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		errorResponse(c, 500, "获取用户列表失败")
		return
	}

	successResponse(c, gin.H{
		"users": users,
		"total": len(users),
	})
}

// 根据ID获取用户
func getUserById(c *gin.Context) {
	id := c.Param("id")
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		errorResponse(c, 404, "用户不存在")
		return
	}

	successResponse(c, user)
}

// 创建用户
func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 检查用户名是否存在
	var existingUser User
	if err := db.Where("username = ? OR email = ?", newUser.Username, newUser.Email).First(&existingUser).Error; err == nil {
		errorResponse(c, 400, "用户名或邮箱已存在")
		return
	}

	// 加密密码
	if newUser.Password != "" {
		hashedPassword, err := hashPassword(newUser.Password)
		if err != nil {
			errorResponse(c, 500, "密码加密失败")
			return
		}
		newUser.Password = hashedPassword
	}

	result := db.Create(&newUser)
	if result.Error != nil {
		errorResponse(c, 500, "创建用户失败")
		return
	}

	// 不返回密码
	newUser.Password = ""
	successResponse(c, newUser)
}

// 更新用户
func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	
	// 查找用户
	result := db.First(&user, id)
	if result.Error != nil {
		errorResponse(c, 404, "用户不存在")
		return
	}

	// 保存原密码
	oldPassword := user.Password

	// 绑定更新数据
	var updateData User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 更新基本字段
	user.Email = updateData.Email
	user.Role = updateData.Role
	user.Status = updateData.Status
	
	// 更新个人资料字段
	user.RealName = updateData.RealName
	user.Phone = updateData.Phone
	user.Avatar = updateData.Avatar
	user.Department = updateData.Department
	user.Position = updateData.Position
	user.Bio = updateData.Bio

	// 处理密码更新
	if updateData.Password != "" {
		hashedPassword, err := hashPassword(updateData.Password)
		if err != nil {
			errorResponse(c, 500, "密码加密失败")
			return
		}
		user.Password = hashedPassword
	} else {
		// 如果没有提供新密码，保持原密码
		user.Password = oldPassword
	}

	// 保存更新
	result = db.Save(&user)
	if result.Error != nil {
		errorResponse(c, 500, "更新用户失败")
		return
	}
	
	// 不返回密码
	user.Password = ""
	successResponse(c, user)
}

// 删除用户
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "删除用户失败")
		return
	}

	if result.RowsAffected == 0 {
		errorResponse(c, 404, "用户不存在")
		return
	}

	successResponse(c, gin.H{"message": "用户删除成功"})
}

// 获取系统信息
func getSystemInfo(c *gin.Context) {
	var userCount int64
	db.Model(&User{}).Count(&userCount)

	successResponse(c, gin.H{
		"user_count":    userCount,
		"database_type": "SQLite",
		"version":       "v1.0.0",
		"uptime":        time.Now().Unix(),
	})
}

// 更新个人资料
func updateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		errorResponse(c, 404, "用户不存在")
		return
	}

	// 绑定更新数据
	var updateData struct {
		Email      string `json:"email"`
		RealName   string `json:"real_name"`
		Phone      string `json:"phone"`
		Avatar     string `json:"avatar"`
		Department string `json:"department"`
		Position   string `json:"position"`
		Bio        string `json:"bio"`
	}
	
	if err := c.ShouldBindJSON(&updateData); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 更新个人资料字段（不允许修改用户名、角色等敏感信息）
	user.Email = updateData.Email
	user.RealName = updateData.RealName
	user.Phone = updateData.Phone
	user.Avatar = updateData.Avatar
	user.Department = updateData.Department
	user.Position = updateData.Position
	user.Bio = updateData.Bio

	// 保存更新
	result = db.Save(&user)
	if result.Error != nil {
		errorResponse(c, 500, "更新个人资料失败")
		return
	}
	
	// 不返回密码
	user.Password = ""
	successResponse(c, user)
}

// 系统监控API
func getSystemMonitor(c *gin.Context) {
	var userCount, activeUserCount, todayLoginCount, logCount, todayLogCount int64
	var dbSize int64

	// 用户总数
	db.Model(&User{}).Count(&userCount)
	// 活跃用户（最近7天登录）
	db.Raw("SELECT COUNT(*) FROM users WHERE last_login >= date('now', '-7 day')").Scan(&activeUserCount)
	// 今日登录用户
	db.Raw("SELECT COUNT(*) FROM users WHERE last_login >= date('now', 'start of day')").Scan(&todayLoginCount)
	// 操作日志总数
	db.Model(&OperationLog{}).Count(&logCount)
	// 今日操作日志
	db.Raw("SELECT COUNT(*) FROM operation_logs WHERE created_at >= date('now', 'start of day')").Scan(&todayLogCount)
	// 数据库文件大小
	if info, err := os.Stat("jing_admin.db"); err == nil {
		dbSize = info.Size()
	}

	successResponse(c, gin.H{
		"user_count": userCount,
		"active_user_count": activeUserCount,
		"today_login_count": todayLoginCount,
		"log_count": logCount,
		"today_log_count": todayLogCount,
		"db_size": dbSize,
		"server_time": time.Now().Format("2006-01-02 15:04:05"),
		"uptime": time.Since(startTime).Seconds(),
	})
}

var startTime = time.Now() 