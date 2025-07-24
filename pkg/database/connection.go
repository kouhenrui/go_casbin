package database

import (
	"fmt"
	"go_casbin/internal/logger"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

type DBOption struct {
	DBType    string  `json:"dbType"`
	DBName    string  `json:"dbName" `
	Username  string  ` json:"username" `
	Password  string  ` json:"password"`
	Host      string  ` json:"host" `
	Port      int     ` json:"port" `
	Charset   string  ` json:"charset" `
	ParseTime *bool   ` json:"parseTime" `
	Loc       *string ` json:"loc" `
}

// InitDB 初始化数据库连接
func InitDB(option DBOption) error {
	var err error
	once.Do(func() {
		db, err = connectDB(option)
		if err != nil {
			logger.ErrorWithErr("数据库连接失败", err)
			return
		}

		// 设置连接池
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)

		logger.Info("数据库连接成功")
	})
	return err
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if db == nil {
		panic("数据库未初始化，请先调用 InitDB()")
	}
	return db
}

// connectDB 连接数据库
func connectDB(option DBOption) (*gorm.DB, error) {
	cfg := option
	if cfg.Loc == nil {
		loc := "Asia/Shanghai"
		cfg.Loc = &loc
	}
	if cfg.ParseTime == nil {
		parseTime := true
		cfg.ParseTime = &parseTime
	}
	// 根据数据库类型创建连接
	switch cfg.DBType {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset, *cfg.ParseTime, *cfg.Loc)
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// 支持postgres
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
			cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port, *cfg.Loc)
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})

	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", cfg.DBType)
	}
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// AutoMigrate 自动迁移表结构
func AutoMigrate(models ...interface{}) error {
	return GetDB().AutoMigrate(models...)
}
