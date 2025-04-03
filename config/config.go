package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	MySQL MySQLConfig
	Redis RedisConfig
	App   AppConfig
}

type MySQLConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}

type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  int
	ReadTimeout  int
	WriteTimeout int
}

type AppConfig struct {
	Host      string
	Port      int
	Mode      string
	JWTSecret string
	JWTExpire int
}

var GlobalConfig Config

func InitConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置文件
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

// GetAppHost 返回应用主机
func (c *AppConfig) GetHost() (string, int) {
	return c.Host, c.Port
}

// GetMySQLDSN 返回MySQL连接字符串
func (c *MySQLConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

// GetJWTKey 返回JWT密钥
func (c *AppConfig) GetJWTKey() string {
	return c.JWTSecret
}
