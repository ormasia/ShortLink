package config

import (
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

// Config 库存服务配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Nacos    NacosConfig    `mapstructure:"nacos"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// NacosConfig Nacos配置
type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	Group     string `mapstructure:"group"`
	DataID    string `mapstructure:"dataid"`
}

// LoadConfig 从Nacos加载配置
func LoadConfig() (*Config, error) {
	clientConfig := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "localhost",
			Port:   8848,
		},
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return nil, err
	}

	// 从本地文件加载默认配置
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("../../config")

	if err := v.ReadInConfig(); err != nil {
		// 如果找不到配置文件，使用默认配置
		v.SetDefault("server.port", 8082)
		v.SetDefault("database.host", "localhost")
		v.SetDefault("database.port", 3306)
		v.SetDefault("database.username", "root")
		v.SetDefault("database.password", "password")
		v.SetDefault("database.dbname", "shortlink")
	}

	// 尝试从Nacos获取配置
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "inventoryservice",
		Group:  "DEFAULT_GROUP",
	})

	if err == nil && content != "" {
		// 解析Nacos配置
		v2 := viper.New()
		v2.SetConfigType("yaml")
		v2.ReadConfig(strings.NewReader(content))

		// 合并配置
		for _, key := range v2.AllKeys() {
			v.Set(key, v2.Get(key))
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
