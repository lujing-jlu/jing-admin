package main

import (
	"encoding/csv"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 导出用户数据为CSV
func exportUsersCSV(c *gin.Context) {
	// 权限检查（已由adminMiddleware保证）
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		errorResponse(c, 500, "获取用户数据失败")
		return
	}

	c.Header("Content-Disposition", "attachment; filename=users_"+time.Now().Format("20060102_150405")+".csv")
	c.Header("Content-Type", "text/csv; charset=utf-8")

	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	// 写入表头
	headers := []string{"ID", "用户名", "邮箱", "角色", "状态", "真实姓名", "手机", "部门", "职位", "简介", "最后登录", "创建时间"}
	w.Write(headers)

	// 写入数据
	for _, user := range users {
		row := []string{
			strconv.FormatUint(uint64(user.ID), 10),
			user.Username,
			user.Email,
			user.Role,
			ifThenElse(user.Status, "启用", "禁用"),
			user.RealName,
			user.Phone,
			user.Department,
			user.Position,
			user.Bio,
			formatTime(user.LastLogin),
			user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		w.Write(row)
	}
}

// 导出角色数据为CSV
func exportRolesCSV(c *gin.Context) {
	var roles []Role
	result := db.Find(&roles)
	if result.Error != nil {
		errorResponse(c, 500, "获取角色数据失败")
		return
	}

	c.Header("Content-Disposition", "attachment; filename=roles_"+time.Now().Format("20060102_150405")+".csv")
	c.Header("Content-Type", "text/csv; charset=utf-8")

	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	// 写入表头
	headers := []string{"ID", "角色名", "显示名", "描述", "状态", "创建时间"}
	w.Write(headers)

	// 写入数据
	for _, role := range roles {
		row := []string{
			strconv.FormatUint(uint64(role.ID), 10),
			role.Name,
			role.DisplayName,
			role.Description,
			ifThenElse(role.Status, "启用", "禁用"),
			role.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		w.Write(row)
	}
}

// 导出权限数据为CSV
func exportPermissionsCSV(c *gin.Context) {
	var permissions []Permission
	result := db.Find(&permissions)
	if result.Error != nil {
		errorResponse(c, 500, "获取权限数据失败")
		return
	}

	c.Header("Content-Disposition", "attachment; filename=permissions_"+time.Now().Format("20060102_150405")+".csv")
	c.Header("Content-Type", "text/csv; charset=utf-8")

	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	// 写入表头
	headers := []string{"ID", "权限名", "显示名", "资源", "操作", "描述", "创建时间"}
	w.Write(headers)

	// 写入数据
	for _, perm := range permissions {
		row := []string{
			strconv.FormatUint(uint64(perm.ID), 10),
			perm.Name,
			perm.DisplayName,
			perm.Resource,
			perm.Action,
			perm.Description,
			perm.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		w.Write(row)
	}
}

// 工具函数：状态转换
func ifThenElse(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

// 工具函数：格式化时间
func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
} 