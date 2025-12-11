package config

import (
	"log"
	"os"
	// "path/filepath"
	"time"
)

// Config 应用配置
type Config struct {
	// 服务器配置
	Server struct {
		Port string
		Mode string // development, production
	}

	// 数据库配置
	Database struct {
		Type     string // mysql, postgres, sqlite
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		Charset  string
		ParseTime bool
		Loc      string
	}

	// JWT配置
	JWT struct {
		Secret    string
		ExpireDur time.Duration
	}

	// 跨域配置
	CORS struct {
		AllowOrigins []string
	}
}

// AppConfig 全局配置实例
var AppConfig Config

// InitConfig 初始化配置
func InitConfig() {
	// 设置默认配置
	setDefaultConfig()

	// 从环境变量加载配置
	loadFromEnv()

	log.Println("配置初始化完成")
}

// 设置默认配置
func setDefaultConfig() {
	// 服务器配置
	AppConfig.Server.Port = "8080"
	AppConfig.Server.Mode = "development"

	// 数据库配置 - 使用你提供的MySQL配置
	AppConfig.Database.Type = "mysql"
	AppConfig.Database.Host = "117.72.38.183"
	AppConfig.Database.Port = "3306"
	AppConfig.Database.User = "root"
	AppConfig.Database.Password = "zyz20050922"
	AppConfig.Database.Name = "gin_vue_chat"
	AppConfig.Database.Charset = "utf8mb4"
	AppConfig.Database.ParseTime = true
	AppConfig.Database.Loc = "Local"

	// JWT配置
	AppConfig.JWT.Secret = "your-secret-key-change-in-production"
	AppConfig.JWT.ExpireDur = 24 * time.Hour

	// CORS配置
	AppConfig.CORS.AllowOrigins = []string{"*"}
}

// 从环境变量加载配置
func loadFromEnv() {
	// 服务器配置
	if port := os.Getenv("SERVER_PORT"); port != "" {
		AppConfig.Server.Port = port
	}
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		AppConfig.Server.Mode = mode
	}

	// 数据库配置
	if dbType := os.Getenv("DB_TYPE"); dbType != "" {
		AppConfig.Database.Type = dbType
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		AppConfig.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		AppConfig.Database.Port = dbPort
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		AppConfig.Database.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		AppConfig.Database.Password = dbPassword
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		AppConfig.Database.Name = dbName
	}
	if dbCharset := os.Getenv("DB_CHARSET"); dbCharset != "" {
		AppConfig.Database.Charset = dbCharset
	}

	// JWT配置
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		AppConfig.JWT.Secret = jwtSecret
	}
}