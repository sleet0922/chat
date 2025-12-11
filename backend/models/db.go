package models

import (
	"fmt"
	"log"
	"time"

	"github.com/yourusername/gin-vue-chat/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局MySQL数据库连接
var DB *gorm.DB

// InitDB 初始化MySQL数据库连接
func InitDB() {
	// 构建MySQL连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Name,
		config.AppConfig.Database.Charset,
		config.AppConfig.Database.ParseTime,
		config.AppConfig.Database.Loc,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("无法连接到MySQL数据库: %v", err)
	}

	// 获取底层的sql.DB对象并设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取数据库连接池失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库表
	err = autoMigrate()
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("成功连接到MySQL数据库")
}

// autoMigrate 自动创建或更新数据库表结构
func autoMigrate() error {
	return DB.AutoMigrate(
		&User{},
		&Friendship{},
		&Group{},
		&GroupMember{},
		&Message{},
	)
}