package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

// 操作日志模型
type OperationLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Username   string    `json:"username" gorm:"not null"`
	Action     string    `json:"action" gorm:"not null"`      // 操作类型：create, update, delete, login, logout
	Resource   string    `json:"resource" gorm:"not null"`    // 操作资源：user, role, permission, system
	ResourceID string    `json:"resource_id"`                 // 资源ID
	Method     string    `json:"method" gorm:"not null"`      // HTTP方法：GET, POST, PUT, DELETE
	Path       string    `json:"path" gorm:"not null"`        // 请求路径
	IP         string    `json:"ip" gorm:"not null"`          // 用户IP
	UserAgent  string    `json:"user_agent"`                  // 用户代理
	Status     int       `json:"status" gorm:"not null"`      // 响应状态码
	Details    string    `json:"details" gorm:"type:text"`    // 详细信息
	CreatedAt  time.Time `json:"created_at"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
}

// 初始化日志系统
func initLogSystem() error {
	// 自动迁移数据库
	err := db.AutoMigrate(&OperationLog{})
	if err != nil {
		return err
	}
	return nil
}

// 记录操作日志
func logOperation(userID uint, username, action, resource, resourceID, method, path, ip, userAgent string, status int, details string) {
	log := OperationLog{
		UserID:     userID,
		Username:   username,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Method:     method,
		Path:       path,
		IP:         ip,
		UserAgent:  userAgent,
		Status:     status,
		Details:    details,
		CreatedAt:  time.Now(),
	}
	
	// 异步写入日志，避免影响主要业务
	go func() {
		db.Create(&log)
	}()
}

// 日志中间件
func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理请求
		c.Next()
		
		// 获取用户信息
		userID, exists := c.Get("user_id")
		if !exists {
			return // 未认证用户不记录日志
		}
		
		username, _ := c.Get("username")
		
		// 确定操作类型和资源
		action := getActionFromMethod(c.Request.Method)
		resource := getResourceFromPath(c.Request.URL.Path)
		resourceID := c.Param("id")
		
		// 获取IP地址
		ip := c.ClientIP()
		
		// 记录操作日志
		if shouldLogOperation(c.Request.URL.Path) {
			uid := userID.(uint)
			uname := username.(string)
			
			logOperation(
				uid,
				uname,
				action,
				resource,
				resourceID,
				c.Request.Method,
				c.Request.URL.Path,
				ip,
				c.Request.UserAgent(),
				c.Writer.Status(),
				"",
			)
		}
	}
}

// 从HTTP方法确定操作类型
func getActionFromMethod(method string) string {
	switch method {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "unknown"
	}
}

// 从请求路径确定资源类型
func getResourceFromPath(path string) string {
	if contains(path, "/users") {
		return "user"
	} else if contains(path, "/roles") {
		return "role"
	} else if contains(path, "/permissions") {
		return "permission"
	} else if contains(path, "/system") {
		return "system"
	} else if contains(path, "/auth") {
		return "auth"
	}
	return "other"
}

// 判断是否需要记录日志
func shouldLogOperation(path string) bool {
	// 不记录的路径
	excludePaths := []string{
		"/health",
		"/api/test",
		"/api/me",
		"/api/my-permissions",
	}
	
	for _, excludePath := range excludePaths {
		if path == excludePath {
			return false
		}
	}
	
	return true
}

// 辅助函数：检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
} 