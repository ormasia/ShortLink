package config

import (
	"fmt"

	"github.com/spf13/viper"
)

//注意，viper识别参数有问题，如果有'_','.'会出现解析不到的情况，使用`mapstructure:"yaml.name"`绑定

type Config struct {
	MySQL  MySQLConfig
	Redis  RedisConfig
	Kafka  KafkaConfig
	Logger LoggerConfig
	App    AppConfig
	Nacos  NacosConfig
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
	Host         string
	Port         int
	Mode         string
	JWTSecret    string
	JWTExpire    int
	Base62Length int
	MaxRetries   int `mapstructure:"max_retries"`
}

type NacosConfig struct {
	ServiceName string
	GroupName   string
	Namespace   string
	Weight      int
	Enabled     bool
	Ip          string
	Port        int
	Metadata    map[string]string
}

var GlobalConfig Config

// InitConfig 初始化配置
// 参数：
//   - configPath: 配置文件路径
//
// 返回：
//   - error: 错误信息
func InitConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// logger.Log.Error("读取配置文件失败", zap.Error(err))
		return err
	}

	// 解析配置文件
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		// logger.Log.Error("解析配置文件失败", zap.Error(err))
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
