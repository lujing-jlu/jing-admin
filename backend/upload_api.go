package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"image"
	"image/jpeg"
	"image/png"
	"github.com/disintegration/imaging"

	"github.com/gin-gonic/gin"
)

// 文件上传模型
type UploadedFile struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OriginalName string    `json:"original_name" gorm:"not null"`    // 原始文件名
	FileName    string    `json:"file_name" gorm:"unique;not null"`  // 存储文件名
	FilePath    string    `json:"file_path" gorm:"not null"`        // 文件路径
	FileSize    int64     `json:"file_size" gorm:"not null"`        // 文件大小(字节)
	FileType    string    `json:"file_type" gorm:"not null"`        // 文件类型(扩展名)
	MimeType    string    `json:"mime_type" gorm:"not null"`        // MIME类型
	UserID      uint      `json:"user_id" gorm:"not null"`          // 上传用户ID
	Username    string    `json:"username" gorm:"not null"`         // 上传用户名
	Category    string    `json:"category" gorm:"default:general"`  // 文件分类: avatar, document, image, general
	IsPublic    bool      `json:"is_public" gorm:"default:false"`   // 是否为公开文件
	CreatedAt   time.Time `json:"created_at"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
}

// 初始化上传系统
func initUploadSystem() error {
	// 自动迁移数据库
	err := db.AutoMigrate(&UploadedFile{})
	if err != nil {
		return err
	}

	// 创建上传目录
	uploadDirs := []string{
		"uploads",
		"uploads/avatars",
		"uploads/documents",
		"uploads/images",
		"uploads/general",
	}

	for _, dir := range uploadDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create upload directory %s: %v", dir, err)
		}
	}

	return nil
}

// 上传文件
func uploadFile(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, 401, "用户未登录")
		return
	}

	username, _ := c.Get("username")

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		errorResponse(c, 400, "请选择要上传的文件")
		return
	}

	// 获取文件分类
	category := c.PostForm("category")
	if category == "" {
		category = "general"
	}

	// 验证文件分类
	validCategories := []string{"avatar", "document", "image", "general"}
	if !containsSlice(validCategories, category) {
		errorResponse(c, 400, "无效的文件分类")
		return
	}

	// 检查文件大小限制
	maxSize := getConfigValue("upload_max_size", "10485760") // 默认10MB
	maxSizeInt, _ := strconv.ParseInt(maxSize, 10, 64)
	if file.Size > maxSizeInt {
		errorResponse(c, 400, fmt.Sprintf("文件大小超出限制(最大 %.1f MB)", float64(maxSizeInt)/(1024*1024)))
		return
	}

	// 检查文件类型
	allowedTypes := getConfigValue("upload_allowed_types", "jpg,jpeg,png,gif,pdf,doc,docx,xls,xlsx")
	fileExt := strings.ToLower(filepath.Ext(file.Filename)[1:])
	if fileExt == "" {
		errorResponse(c, 400, "文件必须有扩展名")
		return
	}

	allowedTypesList := strings.Split(allowedTypes, ",")
	if !containsSlice(allowedTypesList, fileExt) {
		errorResponse(c, 400, fmt.Sprintf("不支持的文件类型: %s", fileExt))
		return
	}

	// 生成唯一文件名
	timestamp := time.Now().Unix()
	newFileName := fmt.Sprintf("%d_%d.%s", userID, timestamp, fileExt)

	// 确定存储路径
	var subDir string
	switch category {
	case "avatar":
		subDir = "avatars"
	case "document":
		subDir = "documents"
	case "image":
		subDir = "images"
	default:
		subDir = "general"
	}

	relativePath := filepath.Join("uploads", subDir, newFileName)
	fullPath := filepath.Join(".", relativePath)

	// 保存文件
	if category == "avatar" {
		// 头像特殊处理：自动裁剪最大正方形并缩放为256x256
		src, err := file.Open()
		if err != nil {
			errorResponse(c, 500, "文件读取失败")
			return
		}
		defer src.Close()
		img, format, err := image.Decode(src)
		if err != nil {
			errorResponse(c, 400, "不支持的图片格式")
			return
		}
		// 裁剪最大正方形
		minSize := img.Bounds().Dx()
		if img.Bounds().Dy() < minSize {
			minSize = img.Bounds().Dy()
		}
		x0 := (img.Bounds().Dx() - minSize) / 2
		y0 := (img.Bounds().Dy() - minSize) / 2
		cropRect := image.Rect(x0, y0, x0+minSize, y0+minSize)
		cropped := imaging.Crop(img, cropRect)
		thumb := imaging.Resize(cropped, 256, 256, imaging.Lanczos)
		out, err := os.Create(fullPath)
		if err != nil {
			errorResponse(c, 500, "头像保存失败")
			return
		}
		defer out.Close()
		switch format {
		case "jpeg":
			jpeg.Encode(out, thumb, &jpeg.Options{Quality: 92})
		case "png":
			png.Encode(out, thumb)
		default:
			jpeg.Encode(out, thumb, &jpeg.Options{Quality: 92})
		}
	} else {
		if err := c.SaveUploadedFile(file, fullPath); err != nil {
			errorResponse(c, 500, "文件保存失败")
			return
		}
	}

	// 获取MIME类型
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = getMimeTypeByExt(fileExt)
	}

	// 检查是否为头像上传（替换旧头像）
	uid := userID.(uint)
	uname := username.(string)
	
	if category == "avatar" {
		// 删除用户之前的头像文件
		var oldAvatars []UploadedFile
		db.Where("user_id = ? AND category = ?", uid, "avatar").Find(&oldAvatars)
		for _, oldAvatar := range oldAvatars {
			// 删除文件
			os.Remove(oldAvatar.FilePath)
			// 删除数据库记录
			db.Delete(&oldAvatar)
		}
	}

	// 保存文件信息到数据库
	uploadedFile := UploadedFile{
		OriginalName: file.Filename,
		FileName:     newFileName,
		FilePath:     fullPath,
		FileSize:     file.Size,
		FileType:     fileExt,
		MimeType:     mimeType,
		UserID:       uid,
		Username:     uname,
		Category:     category,
		IsPublic:     category == "avatar" || category == "image", // 头像和图片默认公开
		CreatedAt:    time.Now(),
	}

	result := db.Create(&uploadedFile)
	if result.Error != nil {
		// 如果数据库保存失败，删除已上传的文件
		os.Remove(fullPath)
		errorResponse(c, 500, "保存文件信息失败")
		return
	}

	// 如果是头像上传，更新用户头像字段
	if category == "avatar" {
		avatarURL := fmt.Sprintf("/api/uploads/%s", newFileName)
		db.Model(&User{}).Where("id = ?", uid).Update("avatar", avatarURL)
	}

	successResponse(c, gin.H{
		"file": uploadedFile,
		"url":  fmt.Sprintf("/api/uploads/%s", newFileName),
	})
}

// 获取文件列表
func getFileList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, 401, "用户未登录")
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")

	offset := (page - 1) * pageSize

	// 构建查询
	query := db.Model(&UploadedFile{})
	
	// 如果不是管理员，只能看到自己的文件
	role, _ := c.Get("role")
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取文件列表
	var files []UploadedFile
	result := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&files)
	if result.Error != nil {
		errorResponse(c, 500, "获取文件列表失败")
		return
	}

	successResponse(c, gin.H{
		"files": files,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// 删除文件
func deleteFile(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, 401, "用户未登录")
		return
	}

	var file UploadedFile
	result := db.First(&file, id)
	if result.Error != nil {
		errorResponse(c, 404, "文件不存在")
		return
	}

	// 检查权限（只能删除自己的文件，除非是管理员）
	role, _ := c.Get("role")
	uid := userID.(uint)
	if role != "admin" && file.UserID != uid {
		errorResponse(c, 403, "没有权限删除此文件")
		return
	}

	// 删除物理文件
	if err := os.Remove(file.FilePath); err != nil {
		// 文件可能已经被删除，记录日志但继续删除数据库记录
		fmt.Printf("Warning: Failed to delete file %s: %v\n", file.FilePath, err)
	}

	// 删除数据库记录
	result = db.Delete(&file)
	if result.Error != nil {
		errorResponse(c, 500, "删除文件记录失败")
		return
	}

	successResponse(c, gin.H{"message": "文件删除成功"})
}

// 下载/访问文件
func serveFile(c *gin.Context) {
	fileName := c.Param("filename")

	// 查找文件信息
	var file UploadedFile
	result := db.Where("file_name = ?", fileName).First(&file)
	if result.Error != nil {
		errorResponse(c, 404, "文件不存在")
		return
	}

	// 检查文件访问权限
	if !file.IsPublic {
		// 私有文件需要登录且只能访问自己的文件
		userID, exists := c.Get("user_id")
		if !exists {
			errorResponse(c, 401, "需要登录才能访问此文件")
			return
		}

		role, _ := c.Get("role")
		uid := userID.(uint)
		if role != "admin" && file.UserID != uid {
			errorResponse(c, 403, "没有权限访问此文件")
			return
		}
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		errorResponse(c, 404, "文件不存在")
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", file.OriginalName))
	c.Header("Content-Type", file.MimeType)

	// 提供文件
	c.File(file.FilePath)
}

// 获取文件统计信息
func getFileStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, 401, "用户未登录")
		return
	}

	role, _ := c.Get("role")
	uid := userID.(uint)

	var stats struct {
		TotalFiles     int64   `json:"total_files"`
		TotalSize      int64   `json:"total_size"`
		AvatarCount    int64   `json:"avatar_count"`
		DocumentCount  int64   `json:"document_count"`
		ImageCount     int64   `json:"image_count"`
		GeneralCount   int64   `json:"general_count"`
	}

	// 构建基础查询
	query := db.Model(&UploadedFile{})
	if role != "admin" {
		query = query.Where("user_id = ?", uid)
	}

	// 总文件数和大小
	query.Count(&stats.TotalFiles)
	query.Select("COALESCE(SUM(file_size), 0)").Scan(&stats.TotalSize)

	// 分类统计
	if role == "admin" {
		db.Model(&UploadedFile{}).Where("category = ?", "avatar").Count(&stats.AvatarCount)
		db.Model(&UploadedFile{}).Where("category = ?", "document").Count(&stats.DocumentCount)
		db.Model(&UploadedFile{}).Where("category = ?", "image").Count(&stats.ImageCount)
		db.Model(&UploadedFile{}).Where("category = ?", "general").Count(&stats.GeneralCount)
	} else {
		db.Model(&UploadedFile{}).Where("user_id = ? AND category = ?", uid, "avatar").Count(&stats.AvatarCount)
		db.Model(&UploadedFile{}).Where("user_id = ? AND category = ?", uid, "document").Count(&stats.DocumentCount)
		db.Model(&UploadedFile{}).Where("user_id = ? AND category = ?", uid, "image").Count(&stats.ImageCount)
		db.Model(&UploadedFile{}).Where("user_id = ? AND category = ?", uid, "general").Count(&stats.GeneralCount)
	}

	successResponse(c, stats)
}

// 辅助函数：根据扩展名获取MIME类型
func getMimeTypeByExt(ext string) string {
	mimeTypes := map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"pdf":  "application/pdf",
		"doc":  "application/msword",
		"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"xls":  "application/vnd.ms-excel",
		"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"txt":  "text/plain",
		"csv":  "text/csv",
	}
	
	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}
	return "application/octet-stream"
}

// 辅助函数：获取系统配置值
func getConfigValue(key, defaultValue string) string {
	var config SystemConfig
	if err := db.Where("key = ?", key).First(&config).Error; err != nil {
		return defaultValue
	}
	return config.Value
}

// 辅助函数：检查字符串切片是否包含某个元素
func containsSlice(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == item {
			return true
		}
	}
	return false
} 