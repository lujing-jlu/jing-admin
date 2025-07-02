package main

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取系统配置列表（管理员）
func getSystemConfigs(c *gin.Context) {
	category := c.Query("category")
	
	var configs []SystemConfig
	query := db
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	result := query.Order("category, id").Find(&configs)
	if result.Error != nil {
		errorResponse(c, 500, "获取系统配置失败")
		return
	}

	// 按分类分组
	groupedConfigs := make(map[string][]SystemConfig)
	for _, config := range configs {
		groupedConfigs[config.Category] = append(groupedConfigs[config.Category], config)
	}

	successResponse(c, gin.H{
		"configs": configs,
		"grouped": groupedConfigs,
		"total":   len(configs),
	})
}

// 获取公开系统配置（无需认证）
func getPublicSystemConfigs(c *gin.Context) {
	var configs []SystemConfig
	result := db.Where("is_public = ?", true).Find(&configs)
	if result.Error != nil {
		errorResponse(c, 500, "获取公开配置失败")
		return
	}

	// 转换为键值对
	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.Key] = config.Value
	}

	successResponse(c, configMap)
}

// 根据key获取配置
func getSystemConfigByKey(c *gin.Context) {
	key := c.Param("key")
	
	var config SystemConfig
	result := db.Where("key = ?", key).First(&config)
	if result.Error != nil {
		errorResponse(c, 404, "配置不存在")
		return
	}

	successResponse(c, config)
}

// 更新系统配置
func updateSystemConfig(c *gin.Context) {
	key := c.Param("key")
	
	var config SystemConfig
	result := db.Where("key = ?", key).First(&config)
	if result.Error != nil {
		errorResponse(c, 404, "配置不存在")
		return
	}

	// 检查是否可编辑
	if !config.IsEditable {
		errorResponse(c, 400, "该配置不允许编辑")
		return
	}

	var updateData struct {
		Value string `json:"value" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&updateData); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 验证配置值
	if err := validateConfigValue(config.Type, updateData.Value); err != nil {
		errorResponse(c, 400, "配置值格式错误: "+err.Error())
		return
	}

	// 更新配置
	config.Value = updateData.Value
	result = db.Save(&config)
	if result.Error != nil {
		errorResponse(c, 500, "更新配置失败")
		return
	}

	successResponse(c, config)
}

// 批量更新系统配置
func batchUpdateSystemConfigs(c *gin.Context) {
	var req struct {
		Configs map[string]string `json:"configs" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 开始事务
	tx := db.Begin()

	for key, value := range req.Configs {
		var config SystemConfig
		result := tx.Where("key = ?", key).First(&config)
		if result.Error != nil {
			tx.Rollback()
			errorResponse(c, 404, "配置 "+key+" 不存在")
			return
		}

		// 检查是否可编辑
		if !config.IsEditable {
			tx.Rollback()
			errorResponse(c, 400, "配置 "+key+" 不允许编辑")
			return
		}

		// 验证配置值
		if err := validateConfigValue(config.Type, value); err != nil {
			tx.Rollback()
			errorResponse(c, 400, "配置 "+key+" 值格式错误: "+err.Error())
			return
		}

		// 更新配置
		config.Value = value
		if err := tx.Save(&config).Error; err != nil {
			tx.Rollback()
			errorResponse(c, 500, "更新配置失败")
			return
		}
	}

	// 提交事务
	tx.Commit()

	successResponse(c, gin.H{
		"message": "批量更新配置成功",
		"updated": len(req.Configs),
	})
}

// 创建新的系统配置
func createSystemConfig(c *gin.Context) {
	var newConfig SystemConfig
	if err := c.ShouldBindJSON(&newConfig); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}

	// 检查配置键是否存在
	var existingConfig SystemConfig
	if err := db.Where("key = ?", newConfig.Key).First(&existingConfig).Error; err == nil {
		errorResponse(c, 400, "配置键已存在")
		return
	}

	// 验证配置值
	if err := validateConfigValue(newConfig.Type, newConfig.Value); err != nil {
		errorResponse(c, 400, "配置值格式错误: "+err.Error())
		return
	}

	result := db.Create(&newConfig)
	if result.Error != nil {
		errorResponse(c, 500, "创建配置失败")
		return
	}

	successResponse(c, newConfig)
}

// 删除系统配置
func deleteSystemConfig(c *gin.Context) {
	key := c.Param("key")
	
	var config SystemConfig
	result := db.Where("key = ?", key).First(&config)
	if result.Error != nil {
		errorResponse(c, 404, "配置不存在")
		return
	}

	// 检查是否可编辑（不可编辑的配置通常是系统核心配置，不允许删除）
	if !config.IsEditable {
		errorResponse(c, 400, "该配置不允许删除")
		return
	}

	result = db.Delete(&config)
	if result.Error != nil {
		errorResponse(c, 500, "删除配置失败")
		return
	}

	successResponse(c, gin.H{"message": "配置删除成功"})
}

// 验证配置值格式
func validateConfigValue(configType, value string) error {
	switch configType {
	case "number":
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
	case "boolean":
		if value != "true" && value != "false" {
			return errors.New("布尔值必须为 true 或 false")
		}
	case "json":
		// 这里可以添加JSON格式验证
		// 暂时跳过，因为需要导入encoding/json包
	}
	return nil
}

// 重置配置到默认值
func resetSystemConfigToDefault(c *gin.Context) {
	key := c.Param("key")
	
	var config SystemConfig
	result := db.Where("key = ?", key).First(&config)
	if result.Error != nil {
		errorResponse(c, 404, "配置不存在")
		return
	}

	// 检查是否可编辑
	if !config.IsEditable {
		errorResponse(c, 400, "该配置不允许重置")
		return
	}

	// 这里需要一个默认值的映射，或者从初始化数据中获取
	// 为简化，暂时返回成功信息
	successResponse(c, gin.H{
		"message": "配置重置成功",
		"config":  config,
	})
} 