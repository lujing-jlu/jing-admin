package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取角色列表
func getRoleList(c *gin.Context) {
	var roles []Role
	result := db.Preload("Permissions").Find(&roles)
	if result.Error != nil {
		errorResponse(c, 500, "获取角色列表失败")
		return
	}

	successResponse(c, gin.H{
		"roles": roles,
		"total": len(roles),
	})
}

// 根据ID获取角色
func getRoleById(c *gin.Context) {
	id := c.Param("id")
	var role Role
	result := db.Preload("Permissions").First(&role, id)
	if result.Error != nil {
		errorResponse(c, 404, "角色不存在")
		return
	}

	successResponse(c, role)
}

// 创建角色
func createRole(c *gin.Context) {
	var newRole Role
	if err := c.ShouldBindJSON(&newRole); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 检查角色名是否存在
	var existingRole Role
	if err := db.Where("name = ?", newRole.Name).First(&existingRole).Error; err == nil {
		errorResponse(c, 400, "角色名已存在")
		return
	}

	result := db.Create(&newRole)
	if result.Error != nil {
		errorResponse(c, 500, "创建角色失败")
		return
	}

	successResponse(c, newRole)
}

// 更新角色
func updateRole(c *gin.Context) {
	id := c.Param("id")
	var role Role
	
	// 查找角色
	result := db.First(&role, id)
	if result.Error != nil {
		errorResponse(c, 404, "角色不存在")
		return
	}

	// 绑定更新数据
	var updateData Role
	if err := c.ShouldBindJSON(&updateData); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 更新字段
	role.DisplayName = updateData.DisplayName
	role.Description = updateData.Description
	role.Status = updateData.Status

	// 保存更新
	db.Save(&role)
	
	successResponse(c, role)
}

// 删除角色
func deleteRole(c *gin.Context) {
	id := c.Param("id")
	
	// 检查是否为系统预定义角色
	var role Role
	if err := db.First(&role, id).Error; err != nil {
		errorResponse(c, 404, "角色不存在")
		return
	}

	// 防止删除系统默认角色
	if role.Name == "super_admin" || role.Name == "admin" || role.Name == "user" {
		errorResponse(c, 400, "不能删除系统默认角色")
		return
	}

	result := db.Delete(&Role{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "删除角色失败")
		return
	}

	if result.RowsAffected == 0 {
		errorResponse(c, 404, "角色不存在")
		return
	}

	successResponse(c, gin.H{"message": "角色删除成功"})
}

// 获取权限列表
func getPermissionList(c *gin.Context) {
	var permissions []Permission
	result := db.Find(&permissions)
	if result.Error != nil {
		errorResponse(c, 500, "获取权限列表失败")
		return
	}

	// 按资源分组
	groupedPermissions := make(map[string][]Permission)
	for _, perm := range permissions {
		groupedPermissions[perm.Resource] = append(groupedPermissions[perm.Resource], perm)
	}

	successResponse(c, gin.H{
		"permissions": permissions,
		"grouped":     groupedPermissions,
		"total":       len(permissions),
	})
}

// 分配角色权限
func assignRolePermissions(c *gin.Context) {
	var req struct {
		RoleID        uint   `json:"role_id" binding:"required"`
		PermissionIDs []uint `json:"permission_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 查找角色
	var role Role
	if err := db.First(&role, req.RoleID).Error; err != nil {
		errorResponse(c, 404, "角色不存在")
		return
	}

	// 查找权限
	var permissions []Permission
	if len(req.PermissionIDs) > 0 {
		if err := db.Where("id IN ?", req.PermissionIDs).Find(&permissions).Error; err != nil {
			errorResponse(c, 400, "权限不存在")
			return
		}
	}

	// 清除原有权限关联
	db.Model(&role).Association("Permissions").Clear()
	
	// 分配新权限
	if len(permissions) > 0 {
		db.Model(&role).Association("Permissions").Append(&permissions)
	}

	// 返回更新后的角色信息
	db.Preload("Permissions").First(&role, req.RoleID)
	successResponse(c, gin.H{
		"message": "权限分配成功",
		"role":    role,
	})
}

// 获取用户权限列表
func getUserPermissionsAPI(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, 401, "用户未登录")
		return
	}

	// 转换userID为uint
	uid, ok := userID.(uint)
	if !ok {
		errorResponse(c, 500, "用户ID类型错误")
		return
	}

	permissions := getUserPermissions(uid)
	
	successResponse(c, gin.H{
		"permissions": permissions,
		"total":       len(permissions),
	})
}

// 分配用户角色
func assignUserRoles(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		errorResponse(c, 400, "用户ID格式错误")
		return
	}

	var req struct {
		RoleIDs []uint `json:"role_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 查找用户
	var user User
	if err := db.First(&user, uint(userID)).Error; err != nil {
		errorResponse(c, 404, "用户不存在")
		return
	}

	// 查找角色
	var roles []Role
	if len(req.RoleIDs) > 0 {
		if err := db.Where("id IN ?", req.RoleIDs).Find(&roles).Error; err != nil {
			errorResponse(c, 400, "角色不存在")
			return
		}
	}

	// 清除原有角色关联
	db.Model(&user).Association("Roles").Clear()
	
	// 分配新角色
	if len(roles) > 0 {
		db.Model(&user).Association("Roles").Append(&roles)
	}

	successResponse(c, gin.H{
		"message": "角色分配成功",
		"user_id": userID,
	})
} 