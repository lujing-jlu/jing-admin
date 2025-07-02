package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// 导入用户数据（CSV）
func importUsersCSV(c *gin.Context) {
	// 权限检查（已由adminMiddleware保证）
	file, err := c.FormFile("file")
	if err != nil {
		errorResponse(c, 400, "请上传CSV文件")
		return
	}

	f, err := file.Open()
	if err != nil {
		errorResponse(c, 500, "文件读取失败")
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1 // 允许不定列数

	headers, err := r.Read()
	if err != nil {
		errorResponse(c, 400, "CSV文件格式错误")
		return
	}

	// 映射表头
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.TrimSpace(h)] = i
	}

	required := []string{"用户名", "邮箱", "角色", "状态"}
	for _, col := range required {
		if _, ok := headerMap[col]; !ok {
			errorResponse(c, 400, fmt.Sprintf("缺少必需列: %s", col))
			return
		}
	}

	// 统计
	successCount := 0
	errorRows := []string{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errorRows = append(errorRows, fmt.Sprintf("读取行失败: %v", err))
			continue
		}

		username := strings.TrimSpace(record[headerMap["用户名"]])
		email := strings.TrimSpace(record[headerMap["邮箱"]])
		role := strings.TrimSpace(record[headerMap["角色"]])
		statusStr := strings.TrimSpace(record[headerMap["状态"]])
		status := statusStr == "启用" || statusStr == "1" || strings.ToLower(statusStr) == "true"

		if username == "" || email == "" || role == "" {
			errorRows = append(errorRows, fmt.Sprintf("用户名/邮箱/角色不能为空: %v", record))
			continue
		}

		// 检查是否已存在
		var existing User
		if err := db.Where("username = ? OR email = ?", username, email).First(&existing).Error; err == nil {
			errorRows = append(errorRows, fmt.Sprintf("用户已存在: %s", username))
			continue
		}

		// 默认密码
		defaultPassword := "123456"
		hashedPassword, _ := hashPassword(defaultPassword)

		user := User{
			Username: username,
			Email: email,
			Password: hashedPassword,
			Role: role,
			Status: status,
		}
		if err := db.Create(&user).Error; err != nil {
			errorRows = append(errorRows, fmt.Sprintf("导入失败: %s (%v)", username, err))
			continue
		}
		successCount++
	}

	successResponse(c, gin.H{
		"success": successCount,
		"failed": len(errorRows),
		"errors": errorRows,
	})
}

// 导入角色数据（CSV）
func importRolesCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		errorResponse(c, 400, "请上传CSV文件")
		return
	}

	f, err := file.Open()
	if err != nil {
		errorResponse(c, 500, "文件读取失败")
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	headers, err := r.Read()
	if err != nil {
		errorResponse(c, 400, "CSV文件格式错误")
		return
	}

	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.TrimSpace(h)] = i
	}

	required := []string{"角色名", "显示名", "状态"}
	for _, col := range required {
		if _, ok := headerMap[col]; !ok {
			errorResponse(c, 400, fmt.Sprintf("缺少必需列: %s", col))
			return
		}
	}

	successCount := 0
	errorRows := []string{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errorRows = append(errorRows, fmt.Sprintf("读取行失败: %v", err))
			continue
		}

		name := strings.TrimSpace(record[headerMap["角色名"]])
		displayName := strings.TrimSpace(record[headerMap["显示名"]])
		statusStr := strings.TrimSpace(record[headerMap["状态"]])
		status := statusStr == "启用" || statusStr == "1" || strings.ToLower(statusStr) == "true"
		desc := ""
		if idx, ok := headerMap["描述"]; ok {
			desc = strings.TrimSpace(record[idx])
		}

		if name == "" || displayName == "" {
			errorRows = append(errorRows, fmt.Sprintf("角色名/显示名不能为空: %v", record))
			continue
		}

		// 检查是否已存在
		var existing Role
		if err := db.Where("name = ?", name).First(&existing).Error; err == nil {
			errorRows = append(errorRows, fmt.Sprintf("角色已存在: %s", name))
			continue
		}

		role := Role{
			Name: name,
			DisplayName: displayName,
			Description: desc,
			Status: status,
		}
		if err := db.Create(&role).Error; err != nil {
			errorRows = append(errorRows, fmt.Sprintf("导入失败: %s (%v)", name, err))
			continue
		}
		successCount++
	}

	successResponse(c, gin.H{
		"success": successCount,
		"failed": len(errorRows),
		"errors": errorRows,
	})
}

// 导入权限数据（CSV）
func importPermissionsCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		errorResponse(c, 400, "请上传CSV文件")
		return
	}

	f, err := file.Open()
	if err != nil {
		errorResponse(c, 500, "文件读取失败")
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	headers, err := r.Read()
	if err != nil {
		errorResponse(c, 400, "CSV文件格式错误")
		return
	}

	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.TrimSpace(h)] = i
	}

	required := []string{"权限名", "显示名", "资源", "操作"}
	for _, col := range required {
		if _, ok := headerMap[col]; !ok {
			errorResponse(c, 400, fmt.Sprintf("缺少必需列: %s", col))
			return
		}
	}

	successCount := 0
	errorRows := []string{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errorRows = append(errorRows, fmt.Sprintf("读取行失败: %v", err))
			continue
		}

		name := strings.TrimSpace(record[headerMap["权限名"]])
		displayName := strings.TrimSpace(record[headerMap["显示名"]])
		resource := strings.TrimSpace(record[headerMap["资源"]])
		action := strings.TrimSpace(record[headerMap["操作"]])
		desc := ""
		if idx, ok := headerMap["描述"]; ok {
			desc = strings.TrimSpace(record[idx])
		}

		if name == "" || displayName == "" || resource == "" || action == "" {
			errorRows = append(errorRows, fmt.Sprintf("必填字段不能为空: %v", record))
			continue
		}

		// 检查是否已存在
		var existing Permission
		if err := db.Where("name = ?", name).First(&existing).Error; err == nil {
			errorRows = append(errorRows, fmt.Sprintf("权限已存在: %s", name))
			continue
		}

		perm := Permission{
			Name: name,
			DisplayName: displayName,
			Resource: resource,
			Action: action,
			Description: desc,
		}
		if err := db.Create(&perm).Error; err != nil {
			errorRows = append(errorRows, fmt.Sprintf("导入失败: %s (%v)", name, err))
			continue
		}
		successCount++
	}

	successResponse(c, gin.H{
		"success": successCount,
		"failed": len(errorRows),
		"errors": errorRows,
	})
} 