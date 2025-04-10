package model

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func InitDB(dataSource string) error {
	var err error
	db, err = gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	// 自动建表
	_ = db.AutoMigrate(&User{})
	return err
}

type URLMapping struct {
	ShortURL    string `gorm:"primaryKey"`
	OriginalURL string `gorm:"not null"`
	UserID      string
	CreateTime  time.Time `gorm:"autoCreateTime"`
}

func (URLMapping) TableName() string {
	return "url_mapping" // 显式指定表名
}

// SaveURLMapping 保存短链接与原始URL的映射关系
// 参数：
//   - shortURL: 短链接
//   - originalURL: 原始URL
//
// 返回：
//   - error: 错误信息，如果保存成功则为nil
func SaveURLMapping(shortURL, originalURL string) error {
	URLMapping := URLMapping{
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
	result := db.Create(&URLMapping)
	return result.Error
}

// GetOriginalURL 获取原始URL
// 参数：
//   - shortURL: 短链接
//
// 返回：
//   - string: 原始URL
//   - error: 错误信息，如果获取成功则为nil
func GetOriginalURL(shortURL string) (string, error) {
	var mapping URLMapping
	result := db.First(&mapping, "short_url = ?", shortURL)
	if result.Error != nil {
		return "", result.Error
	}
	return mapping.OriginalURL, nil
}

// 查询长连接是否存在
// 参数：
//   - originalURL: 原始URL
//
// 返回：
//   - string: 短链接
//   - error: 错误信息，如果查询成功则为nil
func IsOriginalURLExist(originalURL string) string {

	var mapping URLMapping
	db.First(&mapping, "original_url = ?", originalURL)
	return mapping.ShortURL
}

func GetAllShortUrls() []string {
	var urls []string
	// 获取所有短链接
	db.Model(&URLMapping{}).Pluck("short_url", &urls) // pluck是gorm的函数，用于获取单个字段的所有数据
	return urls
}
