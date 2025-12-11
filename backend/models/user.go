package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User MySQL中的用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Email     string    `gorm:"size:100;uniqueIndex" json:"email"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Status    string    `gorm:"size:20;default:'offline'" json:"status"` // online, offline, away
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Friendship MySQL中的好友关系模型
type Friendship struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	FriendID  uint      `gorm:"not null;index" json:"friendId"`
	Status    string    `gorm:"size:20;default:'pending'" json:"status"` // pending, accepted, rejected
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// CreateUser 创建新用户
func CreateUser(username, password, email string) (*User, error) {
	// 检查用户名是否已存在
	var existingUser User
	result := DB.Where("username = ?", username).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("用户名已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 检查邮箱是否已存在
	result = DB.Where("email = ?", email).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("邮箱已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Status:   "offline",
	}

	result = DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	result := DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
func GetUserByID(id uint) (*User, error) {
	var user User
	result := DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(user *User) error {
	result := DB.Save(user)
	return result.Error
}

// AddFriend 添加好友请求
func AddFriend(userID, friendID uint) (*Friendship, error) {
	// 检查用户和好友是否存在
	_, err := GetUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	_, err = GetUserByID(friendID)
	if err != nil {
		return nil, errors.New("好友不存在")
	}

	// 检查是否已经是好友
	var existingFriendship Friendship
	result := DB.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", 
		userID, friendID, friendID, userID).First(&existingFriendship)
	
	if result.Error == nil {
		return nil, errors.New("好友关系已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 创建好友关系
	friendship := &Friendship{
		UserID:   userID,
		FriendID: friendID,
		Status:   "pending",
	}

	result = DB.Create(friendship)
	if result.Error != nil {
		return nil, result.Error
	}

	return friendship, nil
}

// GetFriendships 获取用户的好友关系
func GetFriendships(userID uint, status string) ([]*Friendship, error) {
	var friendships []*Friendship
	query := DB.Where("(user_id = ? OR friend_id = ?)", userID, userID)
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	result := query.Find(&friendships)
	if result.Error != nil {
		return nil, result.Error
	}

	return friendships, nil
}

// UpdateFriendship 更新好友关系
func UpdateFriendship(friendship *Friendship) error {
	result := DB.Save(friendship)
	return result.Error
}

// DeleteFriendship 删除好友关系
func DeleteFriendship(friendshipID uint) error {
	result := DB.Delete(&Friendship{}, friendshipID)
	return result.Error
}