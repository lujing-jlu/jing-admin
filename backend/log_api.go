package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取操作日志列表
func getOperationLogs(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize > 100 {
		pageSize = 100
	}
	
	// 筛选参数
	username := c.Query("username")
	action := c.Query("action")
	resource := c.Query("resource")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	
	// 构建查询
	query := db.Model(&OperationLog{}).Preload("User")
	
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resource != "" {
		query = query.Where("resource = ?", resource)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}
	
	// 获取总数
	var total int64
	query.Count(&total)
	
	// 获取分页数据
	var logs []OperationLog
	offset := (page - 1) * pageSize
	result := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&logs)
	
	if result.Error != nil {
		errorResponse(c, 500, "获取操作日志失败")
		return
	}
	
	successResponse(c, gin.H{
		"logs":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"pages":     (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// 根据ID获取操作日志详情
func getOperationLogById(c *gin.Context) {
	id := c.Param("id")
	var log OperationLog
	result := db.Preload("User").First(&log, id)
	if result.Error != nil {
		errorResponse(c, 404, "操作日志不存在")
		return
	}
	
	successResponse(c, log)
}

// 删除操作日志
func deleteOperationLog(c *gin.Context) {
	id := c.Param("id")
	result := db.Delete(&OperationLog{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "删除操作日志失败")
		return
	}
	
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "操作日志不存在")
		return
	}
	
	successResponse(c, gin.H{"message": "操作日志删除成功"})
}

// 批量删除操作日志
func batchDeleteOperationLogs(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "请求参数错误")
		return
	}
	
	result := db.Where("id IN ?", req.IDs).Delete(&OperationLog{})
	if result.Error != nil {
		errorResponse(c, 500, "批量删除操作日志失败")
		return
	}
	
	successResponse(c, gin.H{
		"message": "批量删除成功",
		"deleted": result.RowsAffected,
	})
}

// 清空操作日志（保留最近30天）
func clearOldOperationLogs(c *gin.Context) {
	// 删除30天前的日志
	result := db.Where("created_at < DATE('now', '-30 day')").Delete(&OperationLog{})
	if result.Error != nil {
		errorResponse(c, 500, "清空旧日志失败")
		return
	}
	
	successResponse(c, gin.H{
		"message": "清空旧日志成功",
		"deleted": result.RowsAffected,
	})
}

// 获取操作日志统计信息
func getOperationLogStats(c *gin.Context) {
	var stats struct {
		TotalLogs      int64                    `json:"total_logs"`
		TodayLogs      int64                    `json:"today_logs"`
		WeekLogs       int64                    `json:"week_logs"`
		ActionStats    []map[string]interface{} `json:"action_stats"`
		ResourceStats  []map[string]interface{} `json:"resource_stats"`
		UserStats      []map[string]interface{} `json:"user_stats"`
	}
	
	// 总日志数
	db.Model(&OperationLog{}).Count(&stats.TotalLogs)
	
	// 今日日志数
	db.Model(&OperationLog{}).Where("DATE(created_at) = DATE('now')").Count(&stats.TodayLogs)
	
	// 本周日志数
	db.Model(&OperationLog{}).Where("created_at >= DATE('now', '-7 day')").Count(&stats.WeekLogs)
	
	// 按操作类型统计
	db.Model(&OperationLog{}).
		Select("action, COUNT(*) as count").
		Group("action").
		Scan(&stats.ActionStats)
	
	// 按资源类型统计
	db.Model(&OperationLog{}).
		Select("resource, COUNT(*) as count").
		Group("resource").
		Scan(&stats.ResourceStats)
	
	// 按用户统计（前10名）
	db.Model(&OperationLog{}).
		Select("username, COUNT(*) as count").
		Group("username").
		Order("count DESC").
		Limit(10).
		Scan(&stats.UserStats)
	
	successResponse(c, stats)
} 