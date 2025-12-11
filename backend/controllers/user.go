package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-vue-chat/models"
	"golang.org/x/crypto/bcrypt"
)

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Avatar string `json:"avatar"`
	Status string `json:"status" binding:"omitempty,oneof=online offline away"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// GetUserProfile 获取用户资料
func GetUserProfile(c *gin.Context) {
	userIDStr := c.GetString("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := models.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"avatar":   user.Avatar,
			"status":   user.Status,
		},
	})
}

// UpdateUserProfile 更新用户资料
func UpdateUserProfile(c *gin.Context) {
	userIDStr := c.GetString("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	user, err := models.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 更新字段
	if req.Email != "" && req.Email != user.Email {
		// 检查邮箱是否已被其他用户使用
		var existingUser models.User
		err = models.DB.Where("email = ? AND id != ?", req.Email, user.ID).First(&existingUser).Error
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "邮箱已被使用"})
			return
		}

		user.Email = req.Email
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if req.Status != "" {
		user.Status = req.Status
	}

	err = models.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户资料失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户资料已更新",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"avatar":   user.Avatar,
			"status":   user.Status,
		},
	})
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	userIDStr := c.GetString("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	user, err := models.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 验证旧密码
	if !user.CheckPassword(req.OldPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "旧密码不正确"})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 更新密码
	user.Password = string(hashedPassword)
	err = models.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码已更新"})
}