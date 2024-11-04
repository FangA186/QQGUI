package sqlconfig

import (
	"QQGUI/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

func MysqlConnect() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:xflovexz@tcp(127.0.0.1:3306)/qq?charset=utf8mb4&parseTime=True&loc=Local",
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get DB instance")
	}
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetMaxIdleConns(50)           // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接的最大生命周期

	if err := db.AutoMigrate(&model.Message{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	return db
}
