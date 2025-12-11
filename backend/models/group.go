package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Group MySQL中的群组模型
type Group struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Avatar      string    `gorm:"size:255" json:"avatar"`
	CreatorID   uint      `gorm:"not null;index" json:"creatorId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `gorm:"index" json:"-"`
}

// GroupMember MySQL中的群组成员模型
type GroupMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GroupID   uint      `gorm:"not null;index" json:"groupId"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	Role      string    `gorm:"size:20;default:'member'" json:"role"` // admin, member
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `gorm:"index" json:"-"`
}

// CreateGroup 创建新群组
func CreateGroup(name, description, avatar string, creatorID uint) (*Group, error) {
	// 检查群组名是否已存在
	var existingGroup Group
	result := DB.Where("name = ?", name).First(&existingGroup)
	if result.Error == nil {
		return nil, errors.New("群组名已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 创建群组
	group := &Group{
		Name:        name,
		Description: description,
		Avatar:      avatar,
		CreatorID:   creatorID,
	}

	result = DB.Create(group)
	if result.Error != nil {
		return nil, result.Error
	}

	// 添加创建者为管理员
	member := &GroupMember{
		GroupID: group.ID,
		UserID:  creatorID,
		Role:    "admin",
	}

	result = DB.Create(member)
	if result.Error != nil {
		return nil, result.Error
	}

	return group, nil
}

// GetGroupByID 根据ID获取群组
func GetGroupByID(id uint) (*Group, error) {
	var group Group
	result := DB.First(&group, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &group, nil
}

// GetGroupsByUserID 获取用户加入的所有群组
func GetGroupsByUserID(userID uint) ([]*Group, error) {
	var groups []*Group
	result := DB.Joins("JOIN group_members ON groups.id = group_members.group_id").
		Where("group_members.user_id = ? AND group_members.deleted_at IS NULL", userID).
		Find(&groups)
	
	if result.Error != nil {
		return nil, result.Error
	}

	return groups, nil
}

// AddGroupMember 添加群组成员
func AddGroupMember(groupID, userID uint, role string) (*GroupMember, error) {
	// 检查群组是否存在
	_, err := GetGroupByID(groupID)
	if err != nil {
		return nil, errors.New("群组不存在")
	}

	// 检查用户是否存在
	_, err = GetUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查用户是否已经是群组成员
	var existingMember GroupMember
	result := DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&existingMember)
	
	if result.Error == nil {
		return nil, errors.New("用户已经是群组成员")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 添加群组成员
	member := &GroupMember{
		GroupID: groupID,
		UserID:  userID,
		Role:    role,
	}

	result = DB.Create(member)
	if result.Error != nil {
		return nil, result.Error
	}

	return member, nil
}

// GetGroupMembers 获取群组成员
func GetGroupMembers(groupID uint) ([]*GroupMember, error) {
	var members []*GroupMember
	result := DB.Where("group_id = ?", groupID).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

// UpdateGroup 更新群组信息
func UpdateGroup(group *Group) error {
	result := DB.Save(group)
	return result.Error
}

// DeleteGroup 删除群组
func DeleteGroup(groupID uint) error {
	result := DB.Delete(&Group{}, groupID)
	return result.Error
}

// RemoveGroupMember 移除群组成员
func RemoveGroupMember(groupID, userID uint) error {
	result := DB.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&GroupMember{})
	return result.Error
}