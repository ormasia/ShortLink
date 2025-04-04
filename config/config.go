package config

import (
	"fmt"
	"shortLink/logger"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	MySQL  MySQLConfig
	Redis  RedisConfig
	Kafka  KafkaConfig
	Logger LoggerConfig
	App    AppConfig
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

type KafkaConfig struct {
	Host    string
	Port    int
	Brokers []string
	Topic   string
}

type LoggerConfig struct {
	Level string
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
		logger.Log.Error("读取配置文件失败", zap.Error(err))
		return err
	}

	// 解析配置文件
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		logger.Log.Error("解析配置文件失败", zap.Error(err))
		return err
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
