package main

import (
	"gorm.io/gorm"
	"time"
)

// 角色模型
type Role struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"unique;not null"`
	DisplayName string       `json:"display_name" gorm:"not null"`
	Description string       `json:"description"`
	Status      bool         `json:"status" gorm:"default:true"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	gorm.Model
}

// 权限模型
type Permission struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"unique;not null"`
	DisplayName string `json:"display_name" gorm:"not null"`
	Resource    string `json:"resource" gorm:"not null"` // 资源名称，如 user, role, system
	Action      string `json:"action" gorm:"not null"`   // 操作名称，如 read, write, delete
	Description string `json:"description"`
	Roles       []Role `json:"roles" gorm:"many2many:role_permissions;"`
	gorm.Model
}

// 用户角色关联模型
type UserRole struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id" gorm:"not null"`
	RoleID uint `json:"role_id" gorm:"not null"`
	User   User `json:"user" gorm:"foreignKey:UserID"`
	Role   Role `json:"role" gorm:"foreignKey:RoleID"`
	gorm.Model
}

// 更新用户模型，支持多角色
type User struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	Username   string     `json:"username" gorm:"unique;not null"`
	Email      string     `json:"email" gorm:"unique;not null"`
	Password   string     `json:"-" gorm:"not null"` // 密码不返回到前端
	Role       string     `json:"role" gorm:"default:user"` // 保留兼容性
	Status     bool       `json:"status" gorm:"default:true"`
	
	// 个人资料字段
	RealName   string     `json:"real_name"`   // 真实姓名
	Phone      string     `json:"phone"`       // 手机号码
	Avatar     string     `json:"avatar"`      // 头像URL
	Department string     `json:"department"`  // 部门
	Position   string     `json:"position"`    // 职位
	Bio        string     `json:"bio"`         // 个人简介
	LastLogin  *time.Time `json:"last_login"`  // 最后登录时间
	
	Roles      []Role     `json:"roles" gorm:"many2many:user_roles;"`
	gorm.Model
}

// 初始化权限系统数据
func initPermissionSystem() error {
	// 自动迁移数据库
	err := db.AutoMigrate(&Role{}, &Permission{}, &UserRole{})
	if err != nil {
		return err
	}

	// 创建默认权限
	defaultPermissions := []Permission{
		{Name: "user.read", DisplayName: "查看用户", Resource: "user", Action: "read", Description: "查看用户列表和详情"},
		{Name: "user.write", DisplayName: "管理用户", Resource: "user", Action: "write", Description: "创建和编辑用户"},
		{Name: "user.delete", DisplayName: "删除用户", Resource: "user", Action: "delete", Description: "删除用户"},
		{Name: "role.read", DisplayName: "查看角色", Resource: "role", Action: "read", Description: "查看角色列表和详情"},
		{Name: "role.write", DisplayName: "管理角色", Resource: "role", Action: "write", Description: "创建和编辑角色"},
		{Name: "role.delete", DisplayName: "删除角色", Resource: "role", Action: "delete", Description: "删除角色"},
		{Name: "permission.read", DisplayName: "查看权限", Resource: "permission", Action: "read", Description: "查看权限配置"},
		{Name: "permission.write", DisplayName: "配置权限", Resource: "permission", Action: "write", Description: "配置角色权限"},
		{Name: "system.read", DisplayName: "查看系统", Resource: "system", Action: "read", Description: "查看系统信息"},
		{Name: "system.write", DisplayName: "管理系统", Resource: "system", Action: "write", Description: "系统配置管理"},
	}

	for _, perm := range defaultPermissions {
		var existingPerm Permission
		if err := db.Where("name = ?", perm.Name).First(&existingPerm).Error; err != nil {
			db.Create(&perm)
		}
	}

	// 创建默认角色
	defaultRoles := []struct {
		Role        Role
		Permissions []string
	}{
		{
			Role: Role{
				Name:        "super_admin",
				DisplayName: "超级管理员",
				Description: "拥有系统所有权限",
				Status:      true,
			},
			Permissions: []string{"user.read", "user.write", "user.delete", "role.read", "role.write", "role.delete", "permission.read", "permission.write", "system.read", "system.write"},
		},
		{
			Role: Role{
				Name:        "admin",
				DisplayName: "管理员",
				Description: "拥有用户管理权限",
				Status:      true,
			},
			Permissions: []string{"user.read", "user.write", "user.delete", "system.read"},
		},
		{
			Role: Role{
				Name:        "user",
				DisplayName: "普通用户",
				Description: "基础用户权限",
				Status:      true,
			},
			Permissions: []string{"user.read"},
		},
	}

	for _, roleData := range defaultRoles {
		var existingRole Role
		if err := db.Where("name = ?", roleData.Role.Name).First(&existingRole).Error; err != nil {
			// 创建角色
			db.Create(&roleData.Role)
			
			// 分配权限
			var permissions []Permission
			db.Where("name IN ?", roleData.Permissions).Find(&permissions)
			db.Model(&roleData.Role).Association("Permissions").Append(&permissions)
		}
	}

	return nil
}

// 检查用户权限
func hasUserPermission(userID uint, permissionName string) bool {
	var count int64
	db.Table("users").
		Select("1").
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("users.id = ? AND permissions.name = ? AND roles.status = true", userID, permissionName).
		Count(&count)
	
	return count > 0
}

// 获取用户所有权限
func getUserPermissions(userID uint) []Permission {
	var permissions []Permission
	db.Table("permissions").
		Select("permissions.*").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN roles ON role_permissions.role_id = roles.id").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.status = true", userID).
		Group("permissions.id").
		Find(&permissions)
	
	return permissions
}

// 系统配置模型
type SystemConfig struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Key         string `json:"key" gorm:"unique;not null"`     // 配置键名
	Value       string `json:"value" gorm:"type:text"`         // 配置值
	Type        string `json:"type" gorm:"default:string"`     // 配置类型：string, number, boolean, json
	Category    string `json:"category" gorm:"not null"`       // 配置分类：basic, mail, security, system
	DisplayName string `json:"display_name" gorm:"not null"`   // 显示名称
	Description string `json:"description"`                    // 配置描述
	IsPublic    bool   `json:"is_public" gorm:"default:false"` // 是否为公开配置（前端可访问）
	IsEditable  bool   `json:"is_editable" gorm:"default:true"` // 是否可编辑
	gorm.Model
}

// 初始化系统配置
func initSystemConfig() error {
	// 自动迁移数据库
	err := db.AutoMigrate(&SystemConfig{})
	if err != nil {
		return err
	}

	// 创建默认系统配置
	defaultConfigs := []SystemConfig{
		// 基础配置
		{Key: "site_name", Value: "Jing Admin", Type: "string", Category: "basic", DisplayName: "网站名称", Description: "系统显示的名称", IsPublic: true, IsEditable: true},
		{Key: "site_description", Value: "现代化管理后台系统", Type: "string", Category: "basic", DisplayName: "网站描述", Description: "系统描述信息", IsPublic: true, IsEditable: true},
		{Key: "site_keywords", Value: "管理系统,后台,Vue3,Go", Type: "string", Category: "basic", DisplayName: "网站关键词", Description: "SEO关键词", IsPublic: true, IsEditable: true},
		{Key: "site_logo", Value: "/logo.png", Type: "string", Category: "basic", DisplayName: "网站Logo", Description: "网站Logo图片地址", IsPublic: true, IsEditable: true},
		{Key: "copyright", Value: "© 2024 Jing Admin. All Rights Reserved.", Type: "string", Category: "basic", DisplayName: "版权信息", Description: "网站版权信息", IsPublic: true, IsEditable: true},
		
		// 邮件配置
		{Key: "mail_host", Value: "", Type: "string", Category: "mail", DisplayName: "邮件服务器", Description: "SMTP服务器地址", IsPublic: false, IsEditable: true},
		{Key: "mail_port", Value: "587", Type: "number", Category: "mail", DisplayName: "邮件端口", Description: "SMTP服务器端口", IsPublic: false, IsEditable: true},
		{Key: "mail_username", Value: "", Type: "string", Category: "mail", DisplayName: "邮箱用户名", Description: "发送邮件的用户名", IsPublic: false, IsEditable: true},
		{Key: "mail_password", Value: "", Type: "string", Category: "mail", DisplayName: "邮箱密码", Description: "发送邮件的密码", IsPublic: false, IsEditable: true},
		{Key: "mail_from", Value: "", Type: "string", Category: "mail", DisplayName: "发件人", Description: "邮件发送者地址", IsPublic: false, IsEditable: true},
		
		// 安全配置
		{Key: "session_timeout", Value: "3600", Type: "number", Category: "security", DisplayName: "会话超时", Description: "用户会话超时时间（秒）", IsPublic: false, IsEditable: true},
		{Key: "password_min_length", Value: "6", Type: "number", Category: "security", DisplayName: "密码最小长度", Description: "用户密码最小长度要求", IsPublic: true, IsEditable: true},
		{Key: "enable_captcha", Value: "false", Type: "boolean", Category: "security", DisplayName: "启用验证码", Description: "登录时是否启用验证码", IsPublic: true, IsEditable: true},
		{Key: "max_login_attempts", Value: "5", Type: "number", Category: "security", DisplayName: "最大登录尝试", Description: "账户锁定前的最大登录尝试次数", IsPublic: false, IsEditable: true},
		
		// 系统配置
		{Key: "system_version", Value: "1.0.0", Type: "string", Category: "system", DisplayName: "系统版本", Description: "当前系统版本号", IsPublic: true, IsEditable: false},
		{Key: "upload_max_size", Value: "10485760", Type: "number", Category: "system", DisplayName: "上传文件大小限制", Description: "文件上传大小限制（字节）", IsPublic: false, IsEditable: true},
		{Key: "upload_allowed_types", Value: "jpg,jpeg,png,gif,pdf,doc,docx,xls,xlsx", Type: "string", Category: "system", DisplayName: "允许上传类型", Description: "允许上传的文件类型", IsPublic: false, IsEditable: true},
		{Key: "pagination_size", Value: "20", Type: "number", Category: "system", DisplayName: "分页大小", Description: "默认分页大小", IsPublic: true, IsEditable: true},
	}

	for _, config := range defaultConfigs {
		var existingConfig SystemConfig
		if err := db.Where("key = ?", config.Key).First(&existingConfig).Error; err != nil {
			db.Create(&config)
		}
	}

	return nil
} 