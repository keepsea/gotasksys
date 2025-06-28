// internal/config/database.go

package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *Config) {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully.")

	// TODO: 在这里自动迁移数据库表结构 (我们将在模型定义后添加)
	// err = DB.AutoMigrate(&model.User{}, &model.Task{}, &model.TaskType{})
	// if err != nil {
	// 	log.Fatalf("Failed to auto migrate database: %v", err)
	// }
}
