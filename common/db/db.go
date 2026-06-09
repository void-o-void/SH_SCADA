package db

import (
	"fmt"
	"SH_SCADA/common/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化 SQLite3 数据库并自动建表
func Init(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印SQL日志，生产环境可改为 logger.Warn
	})
	if err != nil {
		return fmt.Errorf("连接SQLite失败: %w", err)
	}

	// 启用 WAL 模式，提升网关并发读写性能
	DB.Exec("PRAGMA journal_mode=WAL")
	// 设置 busy_timeout，避免写锁冲突
	DB.Exec("PRAGMA busy_timeout=5000")

	// 自动迁移建表
	if err := DB.AutoMigrate(
		&model.DeviceInfo{},
		&model.UnitInfo{},
		&model.ParameterInfo{},
		&model.Task{},
	); err != nil {
		return fmt.Errorf("自动建表失败: %w", err)
	}

	fmt.Println("SQLite3 数据库初始化完成:", dbPath)
	return nil
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
