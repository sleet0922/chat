package models

import (
	"time"
)

// 消息类型常量
const (
	MessageTypePrivate = "private" // 私聊消息
	MessageTypeGroup   = "group"   // 群聊消息
)

// Message MySQL中的消息模型
type Message struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Type       string    `gorm:"size:20;not null" json:"type"` // private, group
	SenderID   uint      `gorm:"not null;index" json:"senderId"`
	ReceiverID uint      `gorm:"index" json:"receiverId,omitempty"` // 私聊时的接收者ID
	GroupID    uint      `gorm:"index" json:"groupId,omitempty"`    // 群聊时的群组ID
	Content    string    `gorm:"type:text;not null" json:"content"`
	Timestamp  time.Time `gorm:"index" json:"timestamp"`
	Read       bool      `gorm:"default:false" json:"read"` // 消息是否已读
}

// SavePrivateMessage 保存私聊消息到MySQL
func SavePrivateMessage(senderID, receiverID uint, content string) (*Message, error) {
	message := &Message{
		Type:       MessageTypePrivate,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Timestamp:  time.Now(),
		Read:       false,
	}

	result := DB.Create(message)
	if result.Error != nil {
		return nil, result.Error
	}

	return message, nil
}

// SaveGroupMessage 保存群聊消息到MySQL
func SaveGroupMessage(senderID, groupID uint, content string) (*Message, error) {
	message := &Message{
		Type:      MessageTypeGroup,
		SenderID:  senderID,
		GroupID:   groupID,
		Content:   content,
		Timestamp: time.Now(),
		Read:      false,
	}

	result := DB.Create(message)
	if result.Error != nil {
		return nil, result.Error
	}

	return message, nil
}

// GetPrivateMessages 获取私聊消息（支持分页）
func GetPrivateMessages(userID, friendID uint, limit, offset int) ([]*Message, error) {
	var messages []*Message
	result := DB.Where("type = ? AND ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))", 
		MessageTypePrivate, userID, friendID, friendID, userID).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages)
	
	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}

// GetGroupMessages 获取群聊消息（支持分页）
func GetGroupMessages(groupID uint, limit, offset int) ([]*Message, error) {
	var messages []*Message
	result := DB.Where("type = ? AND group_id = ?", MessageTypeGroup, groupID).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages)
	
	if result.Error != nil {
		return nil, result.Error
	}

	return messages, nil
}

// MarkMessagesAsRead 标记消息为已读
func MarkMessagesAsRead(messageIDs []uint) error {
	if len(messageIDs) == 0 {
		return nil
	}
	
	result := DB.Model(&Message{}).Where("id IN ?", messageIDs).Update("read", true)
	return result.Error
}