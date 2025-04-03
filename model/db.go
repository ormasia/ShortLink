package model

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db1 *sql.DB

// InitDB 初始化数据库连接
func InitDB1(dsn string) error {
	var err error
	db1, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// 测试连接
	if err := db1.Ping(); err != nil {
		return err
	}

	// 创建必要的表
	if err := createTables(); err != nil {
		return err
	}

	return nil
}

// createTables 创建必要的数据表
func createTables() error {
	// 创建ID生成器表
	_, err := db1.Exec(`
	CREATE TABLE IF NOT EXISTS id_generator (
		id BIGINT AUTO_INCREMENT PRIMARY KEY
	)`)
	if err != nil {
		return err
	}

	// 创建URL映射表
	_, err = db1.Exec(`
	CREATE TABLE IF NOT EXISTS url_mapping (
		short_url VARCHAR(255) PRIMARY KEY,
		original_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}
