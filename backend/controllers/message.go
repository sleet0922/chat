package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gin-vue-chat/models"
	"github.com/yourusername/gin-vue-chat/websocket"
)

// SendPrivateMessageRequest 发送私聊消息请求
type SendPrivateMessageRequest struct {
	ReceiverID string `json:"receiverId" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// SendGroupMessageRequest 发送群聊消息请求
type SendGroupMessageRequest struct {
	GroupID string `json:"groupId" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// GetPrivateMessages 获取私聊消息
func GetPrivateMessages(c *gin.Context) {
	userIDStr := c.GetString("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	receiverIDStr := c.Param("userId")
	receiverID, err := strconv.ParseUint(receiverIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的接收者ID"})
		return
	}

	// 检查接收者是否存在
	_, err = models.GetUserByID(uint(receiverID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 检查是否是好友关系
	friendships, err := models.GetFriendships(uint(userID), "accepted")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 验证是否为好友
	isFriend := false
	for _, friendship := range friendships {
		if (friendship.UserID == uint(userID) && friendship.FriendID == uint(receiverID)) ||
			(friendship.UserID == uint(receiverID) && friendship.FriendID == uint(userID)) {
			isFriend = true
			break
		}
	}

	// 如果不是好友，返回错误
	if !isFriend {
		c.JSON(http.StatusForbidden, gin.H{"error": "您不是该用户的好友"})
		return
	}

	// 获取分页参数
	limit := 20 // 默认每页20条
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0 // 默认从第一条开始
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// 获取消息
	messages, err := models.GetPrivateMessages(uint(userID), uint(receiverID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取消息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// SendPrivateMessage 发送私聊消息
func SendPrivateMessage(c *gin.Context) {
	senderIDStr := c.GetString("userId")
	senderID, err := strconv.ParseUint(senderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req SendPrivateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("请求参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 检查请求参数
	if req.ReceiverID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "接收者ID不能为空"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "消息内容不能为空"})
		return
	}

	// 转换接收者ID
	receiverID, err := strconv.ParseUint(req.ReceiverID, 10, 32)
	if err != nil {
		log.Printf("接收者ID转换失败: %s, 错误: %v", req.ReceiverID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的接收者ID"})
		return
	}

	// 检查接收者是否存在
	_, err = models.GetUserByID(uint(receiverID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "接收者不存在"})
		return
	}

	// 检查是否是好友关系
	friendships, err := models.GetFriendships(uint(senderID), "accepted")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 验证是否为好友
	isFriend := false
	for _, friendship := range friendships {
		if (friendship.UserID == uint(senderID) && friendship.FriendID == uint(receiverID)) ||
			(friendship.UserID == uint(receiverID) && friendship.FriendID == uint(senderID)) {
			isFriend = true
			break
		}
	}

	if !isFriend {
		c.JSON(http.StatusForbidden, gin.H{"error": "您不是该用户的好友"})
		return
	}

	// 保存消息到MySQL
	message, err := models.SavePrivateMessage(uint(senderID), uint(receiverID), req.Content)
	if err != nil {
		log.Printf("保存消息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存消息失败"})
		return
	}

	// 获取发送者信息
	sender, _ := models.GetUserByID(uint(senderID))

	// 通过WebSocket发送消息给接收者
	wsMessage := map[string]interface{}{
		"type": "private",
		"message": map[string]interface{}{
			"id":        message.ID,
			"from":      senderID,
			"to":        receiverID,
			"content":   req.Content,
			"timestamp": message.Timestamp,
			"sender": map[string]interface{}{
				"id":       sender.ID,
				"username": sender.Username,
				"avatar":   sender.Avatar,
			},
		},
	}

	// 获取WebSocket Hub
	hub := c.MustGet("wsHub").(*websocket.Hub)

	// 发送消息给接收者
	// 将消息转换为JSON字符串，再转换为字节数组
	jsonData, err := json.Marshal(gin.H{"data": wsMessage})
	if err != nil {
		// 记录错误但继续执行，因为这不是致命错误
		log.Printf("消息序列化失败: %v", err)
		return
	}
	hub.SendToUser(strconv.FormatUint(uint64(receiverID), 10), jsonData)

	c.JSON(http.StatusOK, gin.H{
		"message": "消息发送成功",
		"data":    message,
	})
}

// GetGroupMessages 获取群聊消息
func GetGroupMessages(c *gin.Context) {
	userIDStr := c.GetString("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	groupIDStr := c.Param("groupId")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的群组ID"})
		return
	}

	// 检查群组是否存在
	_, err = models.GetGroupByID(uint(groupID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "群组不存在"})
		return
	}

	// 检查用户是否是群组成员
	members, err := models.GetGroupMembers(uint(groupID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 验证用户是否在群组中
	isMember := false
	for _, member := range members {
		if member.UserID == uint(userID) {
			isMember = true
			break
		}
	}

	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "您不是该群组的成员"})
		return
	}

	// 获取分页参数
	limit := 20 // 默认每页20条
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0 // 默认从第一条开始
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// 获取消息
	messages, err := models.GetGroupMessages(uint(groupID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取消息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// SendGroupMessage 发送群聊消息
func SendGroupMessage(c *gin.Context) {
	senderIDStr := c.GetString("userId")
	senderID, err := strconv.ParseUint(senderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req SendGroupMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 转换群组ID
	groupID, err := strconv.ParseUint(req.GroupID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的群组ID"})
		return
	}

	// 检查群组是否存在
	_, err = models.GetGroupByID(uint(groupID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "群组不存在"})
		return
	}

	// 检查用户是否是群组成员
	members, err := models.GetGroupMembers(uint(groupID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	// 验证用户是否在群组中
	isMember := false
	for _, member := range members {
		if member.UserID == uint(senderID) {
			isMember = true
			break
		}
	}

	if !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "您不是该群组的成员"})
		return
	}

	// 保存消息到MySQL
	message, err := models.SaveGroupMessage(uint(senderID), uint(groupID), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存消息失败"})
		return
	}

	// 获取发送者信息
	sender, _ := models.GetUserByID(uint(senderID))

	// 通过WebSocket发送消息给群组所有成员
	wsMessage := map[string]interface{}{
		"type": "group",
		"message": map[string]interface{}{
			"id":        message.ID,
			"groupId":   groupID,
			"senderId":  senderID,
			"content":   req.Content,
			"timestamp": message.Timestamp,
			"sender": map[string]interface{}{
				"id":       sender.ID,
				"username": sender.Username,
				"avatar":   sender.Avatar,
			},
		},
	}

	// 获取WebSocket Hub
	hub := c.MustGet("wsHub").(*websocket.Hub)

	// 发送消息给所有成员
	for _, member := range members {
		if member.UserID != uint(senderID) { // 不需要发送给自己
			jsonData, err := json.Marshal(gin.H{"data": wsMessage})
			if err != nil {
				log.Printf("消息序列化失败: %v", err)
				continue
			}
			hub.SendToUser(strconv.FormatUint(uint64(member.UserID), 10), jsonData)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "消息发送成功",
		"data":    message,
	})
}

// MarkMessagesAsRead 标记消息为已读
func MarkMessagesAsRead(c *gin.Context) {
	// 获取消息ID列表
	var req struct {
		MessageIDs []uint `json:"messageIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	if len(req.MessageIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息ID"})
		return
	}

	// 标记消息为已读
	err := models.MarkMessagesAsRead(req.MessageIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "标记消息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "消息已标记为已读"})
}