package model

import "errors"

type IDGenerator struct {
	ID int64 `gorm:"primaryKey"`
}

func (IDGenerator) TableName() string {
	return "id_generator" // 显式指定表名
}

// GenerateID 使用数据库自增主键生成ID
func GenerateID() (int64, error) {
	if db == nil {
		return 0, errors.New("数据库未初始化")
	}

	// 创建新记录
	generator := IDGenerator{}
	//插入到指定表
	result := db.Create(&generator)
	if result.Error != nil {
		return 0, result.Error
	}

	return generator.ID, nil
}

/*
CREATE TABLE id_generator (
    id BIGINT AUTO_INCREMENT PRIMARY KEY
);

CREATE TABLE url_mapping (
    short_url VARCHAR(255) PRIMARY KEY,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
*/
