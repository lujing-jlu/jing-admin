package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT密钥 - 实际项目中应该从环境变量获取
var jwtSecret = []byte("jing-admin-secret-key-2024")

// JWT Claims 结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role,omitempty"`
}

// 登录响应结构
type LoginResponse struct {
	Token    string `json:"token"`
	UserInfo User   `json:"user_info"`
}

// 密码加密
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 密码验证
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 生成JWT令牌
func generateToken(user User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "jing-admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 解析JWT令牌
func parseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// JWT认证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errorResponse(c, 401, "未提供认证令牌")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			errorResponse(c, 401, "认证令牌格式错误")
			c.Abort()
			return
		}

		// 提取令牌
		tokenString := authHeader[7:]
		claims, err := parseToken(tokenString)
		if err != nil {
			errorResponse(c, 401, "无效的认证令牌")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// 管理员权限中间件
func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			errorResponse(c, 403, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

// 登录处理
func login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 查找用户
	var user User
	result := db.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		errorResponse(c, 401, "用户名或密码错误")
		return
	}

	// 验证密码
	if !checkPassword(req.Password, user.Password) {
		errorResponse(c, 401, "用户名或密码错误")
		return
	}

	// 检查用户状态
	if !user.Status {
		errorResponse(c, 401, "账户已被禁用")
		return
	}

	// 生成JWT令牌
	token, err := generateToken(user)
	if err != nil {
		errorResponse(c, 500, "生成令牌失败")
		return
	}

	// 更新最后登录时间
	now := time.Now()
	db.Model(&user).Updates(map[string]interface{}{
		"last_login": &now,
		"updated_at": now,
	})

	// 返回响应（不包含密码）
	user.Password = ""
	successResponse(c, LoginResponse{
		Token:    token,
		UserInfo: user,
	})
}

// 注册处理
func register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 检查用户名是否存在
	var existingUser User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		errorResponse(c, 400, "用户名或邮箱已存在")
		return
	}

	// 加密密码
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		errorResponse(c, 500, "密码加密失败")
		return
	}

	// 设置默认角色
	if req.Role == "" {
		req.Role = "user"
	}

	// 创建新用户
	newUser := User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		Status:   true,
	}

	result := db.Create(&newUser)
	if result.Error != nil {
		errorResponse(c, 500, "创建用户失败")
		return
	}

	// 生成JWT令牌
	token, err := generateToken(newUser)
	if err != nil {
		errorResponse(c, 500, "生成令牌失败")
		return
	}

	// 返回响应（不包含密码）
	newUser.Password = ""
	successResponse(c, LoginResponse{
		Token:    token,
		UserInfo: newUser,
	})
}

// 登出处理（前端删除token即可）
func logout(c *gin.Context) {
	successResponse(c, gin.H{
		"message": "登出成功",
	})
}

// 获取当前用户信息
func getCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		errorResponse(c, 404, "用户不存在")
		return
	}

	// 不返回密码
	user.Password = ""
	successResponse(c, user)
}

// 修改密码
func changePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 获取当前用户
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		errorResponse(c, 404, "用户不存在")
		return
	}

	// 验证旧密码
	if !checkPassword(req.CurrentPassword, user.Password) {
		errorResponse(c, 400, "原密码错误")
		return
	}

	// 加密新密码
	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		errorResponse(c, 500, "密码加密失败")
		return
	}

	// 更新密码
	db.Model(&user).Update("password", hashedPassword)
	
	successResponse(c, gin.H{
		"message": "密码修改成功",
	})
} 