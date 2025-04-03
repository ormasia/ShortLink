package model

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB(dataSource string) error {
	var err error
	db, err = sql.Open("mysql", dataSource)
	return err
}

// SaveURLMapping 保存短链接与原始URL的映射关系
// 参数：
//   - shortURL: 短链接
//   - originalURL: 原始URL
//
// 返回：
//   - error: 错误信息，如果保存成功则为nil
func SaveURLMapping(shortURL, originalURL string) error {
	_, err := db.Exec("INSERT INTO url_mapping (short_url, original_url) VALUES (?, ?)", shortURL, originalURL)
	return err
}

// GetOriginalURL 获取原始URL
// 参数：
//   - shortURL: 短链接
//
// 返回：
//   - string: 原始URL
//   - error: 错误信息，如果获取成功则为nil
func GetOriginalURL(shortURL string) (string, error) {
	var original string
	err := db.QueryRow("SELECT original_url FROM url_mapping WHERE short_url = ?", shortURL).Scan(&original)
	return original, err
}
